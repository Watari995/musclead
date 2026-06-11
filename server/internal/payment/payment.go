// Package payment is the public facade of the payment module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
//
// 設計: ADR 0013 (purchase / payment 分離) + ADR 0014 (Webhook 同期処理)
// 依存方向: purchase → payment (本 module は purchase を知らない)
package payment

import (
	"net/http"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymenthandler "github.com/Watari995/musclead/internal/payment/internal/handler"
	paymentinfra "github.com/Watari995/musclead/internal/payment/internal/infra"
	paymentusecase "github.com/Watari995/musclead/internal/payment/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/go-gorp/gorp/v3"
)

// Module は payment module の公開 API。
// 他 module (purchase 等) は Module.Command / Query 経由でのみ payment を操作できる。
type Module struct {
	Handler http.Handler
	command publicfunctions.PaymentCommand
	query   publicfunctions.PaymentQuery
}

// Command は他 module 公開用 getter (immutable 保護)。
func (m *Module) Command() publicfunctions.PaymentCommand { return m.command }

// Query は他 module 公開用 getter。
func (m *Module) Query() publicfunctions.PaymentQuery { return m.query }

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
	txManager := dbtx.NewTransactionManager(dbmap)

	dbmap.AddTableWithName(paymentinfra.PaymentModel{}, "payments").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.PaymentEventModel{}, "payment_events").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.StripeEventModel{}, "stripe_events").SetKeys(false, "ID")
	dbmap.AddTableWithName(paymentinfra.OutboxEventModel{}, "outbox_events").SetKeys(false, "ID")

	paymentRepo := paymentinfra.NewPaymentRepository(dbmap)
	paymentEventRepo := paymentinfra.NewPaymentEventRepository(dbmap)
	stripeEventRepo := paymentinfra.NewStripeEventRepository(dbmap)
	outboxEventRepo := paymentinfra.NewOutboxEventRepository(dbmap)

	stripeClient := paymentinfra.NewStripeClient(
		cfg.StripeAPIKey,
		cfg.StripeSuccessURL,
		cfg.StripeCancelURL,
		cfg.StripeWebhookSigningSecret,
		cfg.StripePortalReturnURL,
	)
	retryStrategy := &paymentinfra.ExternalRetryStrategy{}

	parseWebhookEvent := paymentusecase.NewParseWebhookEvent(stripeClient)
	initiatePayment := paymentusecase.NewInitiatePayment(paymentRepo, paymentEventRepo, stripeClient)
	completePayment := paymentusecase.NewCompletePayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
	cancelPayment := paymentusecase.NewCancelPayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
	renewPayment := paymentusecase.NewRenewPayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
	handleFailure := paymentusecase.NewHandleFailure(retryStrategy)

	webhookHandler := paymenthandler.NewWebhookHandler(
		parseWebhookEvent,
		completePayment,
		cancelPayment,
		renewPayment,
		handleFailure,
	)

	return &Module{
		Handler: webhookHandler,
		command: initiatePayment,
		query:   paymentQuery{},
	}
}

// paymentQuery は publicfunctions.PaymentQuery の MVP 空実装。 将来 method 追加時に struct を埋める。
type paymentQuery struct{}
