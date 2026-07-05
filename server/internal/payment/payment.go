package payment

import (
	"context"
	"log/slog"
	"time"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	paymentinfra "github.com/Watari995/musclead/internal/payment/internal/infra"
	paymentusecase "github.com/Watari995/musclead/internal/payment/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	outboxinfra "github.com/Watari995/musclead/internal/shared/infra/outbox"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-gorp/gorp/v3"
)

// relayInterval は outbox relay worker のポーリング間隔。
const relayInterval = 10 * time.Second

// Module は payment module の公開 API。
// 他 module (purchase 等) は Module.Command / Query 経由でのみ payment を操作できる。
type Module struct {
	command        publicfunctions.PaymentCommand
	webhookCommand publicfunctions.PaymentWebhookCommand
	query          publicfunctions.PaymentQuery
	processor      publicfunctions.StripeWebhookProcessor
	relayOutbox    *paymentusecase.RelayOutbox
	relayEnabled   bool // SQSQueueURL が設定されている時だけ relay を回す
}

// Config は NewModule に渡す環境差分のまとまり。
type Config struct {
	StripeAPIKey               string
	StripeSuccessURL           string
	StripeCancelURL            string
	StripeWebhookSigningSecret string
	StripePortalReturnURL      string
	SQSQueueURL                string // (旧) SQS relay の送信先。 ResendAPIKey 未設定時のフォールバック
	ResendAPIKey               string // 設定時は Resend で直接メール送信 (SQS/Lambda/SES 不要)
	MailFromAddress            string // メール送信元 (例: no-reply@musclead.com)
}

// NewModule は payment module を初期化する。 Composition Root (cmd/server/main.go) から呼ぶ。
func NewModule(dbmap *gorp.DbMap, cfg Config, userQuery userpublicfunctions.UserQuery, sqsClient *sqs.Client) *Module {
	dbmap.AddTableWithName(paymentinfra.PaymentModel{}, "payments").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.PaymentEventModel{}, "payment_events").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.StripeEventModel{}, "stripe_events").SetKeys(false, "ID")
	dbmap.AddTableWithName(outboxinfra.OutboxEventModel{}, "outbox_events").SetKeys(false, "ID")

	paymentRepo := paymentinfra.NewPaymentRepository(dbmap)
	paymentEventRepo := paymentinfra.NewPaymentEventRepository(dbmap)
	stripeEventRepo := paymentinfra.NewStripeEventRepository(dbmap)
	outboxEventRepo := outboxinfra.NewOutboxEventRepository(dbmap)
	retryStrategy := &paymentinfra.ExternalRetryStrategy{}
	txManager := dbtx.NewTransactionManager(dbmap)

	stripeClient := paymentinfra.NewStripeClient(
		cfg.StripeAPIKey,
		cfg.StripeSuccessURL,
		cfg.StripeCancelURL,
		cfg.StripeWebhookSigningSecret,
		cfg.StripePortalReturnURL,
	)

	initiatePayment := paymentusecase.NewInitiatePayment(paymentRepo, paymentEventRepo, stripeClient)
	createPortalSession := paymentusecase.NewCreatePortalSession(paymentRepo, stripeClient)
	paymentCommand := paymentusecase.NewPaymentCommand(initiatePayment, createPortalSession)

	completePayment := paymentusecase.NewCompletePayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, stripeClient, txManager)
	cancelPayment := paymentusecase.NewCancelPayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
	renewPayment := paymentusecase.NewRenewPayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, stripeClient, txManager)
	handleFailure := paymentusecase.NewHandleFailure(retryStrategy)
	webhookCommand := paymentusecase.NewWebhookCommand(completePayment, cancelPayment, renewPayment, handleFailure)

	parseWebhookEvent := paymentusecase.NewParseWebhookEvent(stripeClient)

	// メール配信経路: Resend (直接送信) を優先。 未設定なら旧 SQS relay にフォールバック。
	// どちらも未設定 (ローカル等) なら relay 無効。
	var publisher paymentdomain.Publisher
	relayEnabled := true
	switch {
	case cfg.ResendAPIKey != "":
		publisher = paymentinfra.NewEmailPublisher(paymentinfra.NewResendMailer(cfg.ResendAPIKey, cfg.MailFromAddress))
	case cfg.SQSQueueURL != "":
		publisher = paymentinfra.NewSQSPublisher(sqsClient, cfg.SQSQueueURL)
	default:
		relayEnabled = false
	}
	relayOutbox := paymentusecase.NewRelayOutbox(outboxEventRepo, paymentRepo, userQuery, publisher)

	return &Module{
		command:        paymentCommand,
		webhookCommand: webhookCommand,
		query:          paymentQuery{}, // 今後追加する場合はここに追記
		processor:      parseWebhookEvent,
		relayOutbox:    relayOutbox,
		relayEnabled:   relayEnabled,
	}
}

// RunRelay は outbox relay worker を起動する。 main.go から goroutine で起動する想定。
// SQSQueueURL 未設定時は no-op (ローカル等で SQS 無しでも起動できるように)。
func (m *Module) RunRelay(ctx context.Context) {
	if !m.relayEnabled {
		slog.Info("outbox relay disabled (no SQS queue URL)")
		return
	}
	ticker := time.NewTicker(relayInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.relayOutbox.Execute(ctx); err != nil {
				slog.Error("outbox relay failed", "err", err)
			}
		}
	}
}

// paymentQuery は publicfunctions.PaymentQuery の MVP 空実装。 将来 method 追加時に struct を埋める。
type paymentQuery struct{}

// Command は purchase 公開用 getter (申込開始)。
func (m *Module) Command() publicfunctions.PaymentCommand { return m.command }

// WebhookCommand は billing 公開用 getter (Webhook 起点の状態遷移)。
func (m *Module) WebhookCommand() publicfunctions.PaymentWebhookCommand { return m.webhookCommand }

// Query は他 module 公開用 getter。
func (m *Module) Query() publicfunctions.PaymentQuery { return m.query }

// Processor は billing 公開用 getter (Stripe Webhook 署名検証 + パース)。
func (m *Module) Processor() publicfunctions.StripeWebhookProcessor { return m.processor }
