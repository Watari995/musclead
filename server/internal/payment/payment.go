// Package payment is the public facade of the payment module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
//
// 設計: ADR 0013 (purchase / payment 分離) + ADR 0014 (Webhook 同期処理)
// 依存方向: purchase → payment (本 module は purchase を知らない)
package payment

import (
	"net/http"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
)

// Module は payment module の公開 API。
// 他 module (purchase 等) は Module.Command / Query 経由でのみ payment を操作できる。
type Module struct {
	// Handler は payment 専用の HTTP handler (例: POST /payment/webhook)。
	// main.go の mux に登録する想定。
	Handler http.Handler

	// command は paymentCommand interface (publicfunctions.PaymentCommand) の実装。
	// 他 module (purchase 等) から ID で参照させるため、 Command() getter 経由でアクセス。
	command publicfunctions.PaymentCommand

	// query は payment の読み込み系 API (publicfunctions.PaymentQuery) の実装。
	query publicfunctions.PaymentQuery
}

// Command は他 module 公開用 getter (immutable 保護)。
func (m *Module) Command() publicfunctions.PaymentCommand { return m.command }

// Query は他 module 公開用 getter。
func (m *Module) Query() publicfunctions.PaymentQuery { return m.query }

// NewModule は payment module を初期化する。 Composition Root (cmd/server/main.go) から呼ぶ。
//
// 必要な引数:
//   - dbmap: gorp DB マッピング
//   - txManager: dbtx.TransactionManager
//   - stripeClient: paymentinfra.NewStripeClient(...)
//   - retryStrategy: paymentinfra.NewExternalRetryStrategy()
//
// TODO (User 実装):
//  1. gorp model 登録 (dbmap.AddTableWithName で PaymentModel / PaymentEventModel / StripeEventModel / OutboxEventModel)
//  2. Repository × 4 を NewXxxRepository で生成
//  3. usecase × 6 を生成:
//     - parseWebhookEvent := paymentusecase.NewParseWebhookEvent(stripeClient)
//     - initiatePayment := paymentusecase.NewInitiatePayment(paymentRepo, paymentEventRepo, stripeClient)
//     - completePayment := paymentusecase.NewCompletePayment(paymentRepo, paymentEventRepo, stripeEventRepo, outboxEventRepo, txManager)
//     - cancelPayment := paymentusecase.NewCancelPayment(...)
//     - renewPayment := paymentusecase.NewRenewPayment(...)
//     - handleFailure := paymentusecase.NewHandleFailure(retryStrategy)
//  4. WebhookHandler を生成 (mux に登録)
//  5. publicfunctions の PaymentCommand / PaymentQuery 実装 struct を作って詰める
//  6. Module を返却
//
// 参考: internal/weight/weight.go, internal/user/user.go
func NewModule( /* TODO: 依存を受け取る引数 */ ) *Module {
	// TODO: User 実装、 上記の wire を組み立てる
	return &Module{}
}
