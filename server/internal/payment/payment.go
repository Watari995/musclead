package payment

import (
	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentinfra "github.com/Watari995/musclead/internal/payment/internal/infra"
	paymentusecase "github.com/Watari995/musclead/internal/payment/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/go-gorp/gorp/v3"
)

// Module は payment module の公開 API。
// 他 module (purchase 等) は Module.Command / Query 経由でのみ payment を操作できる。
type Module struct {
	command        publicfunctions.PaymentCommand
	webhookCommand publicfunctions.PaymentWebhookCommand
	query          publicfunctions.PaymentQuery
	processor      publicfunctions.StripeWebhookProcessor
}

// Config は NewModule に渡す環境差分のまとまり。
type Config struct {
	StripeAPIKey               string
	StripeSuccessURL           string
	StripeCancelURL            string
	StripeWebhookSigningSecret string
	StripePortalReturnURL      string
}

// NewModule は payment module を初期化する。 Composition Root (cmd/server/main.go) から呼ぶ。
func NewModule(dbmap *gorp.DbMap, cfg Config) *Module {
	dbmap.AddTableWithName(paymentinfra.PaymentModel{}, "payments").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.PaymentEventModel{}, "payment_events").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.StripeEventModel{}, "stripe_events").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.OutboxEventModel{}, "outbox_events").SetKeys(false, "ID")

	paymentRepo := paymentinfra.NewPaymentRepository(dbmap)
	paymentEventRepo := paymentinfra.NewPaymentEventRepository(dbmap)
	stripeEventRepo := paymentinfra.NewStripeEventRepository(dbmap)
	outboxEventRepo := paymentinfra.NewOutboxEventRepository(dbmap)
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

	completePayment := paymentusecase.NewCompletePayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
	cancelPayment := paymentusecase.NewCancelPayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
	renewPayment := paymentusecase.NewRenewPayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
	handleFailure := paymentusecase.NewHandleFailure(retryStrategy)
	webhookCommand := paymentusecase.NewWebhookCommand(completePayment, cancelPayment, renewPayment, handleFailure)

	parseWebhookEvent := paymentusecase.NewParseWebhookEvent(stripeClient)

	return &Module{
		command:        initiatePayment,
		webhookCommand: webhookCommand,
		query:          paymentQuery{}, // 今後追加する場合はここに追記
		processor:      parseWebhookEvent,
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
