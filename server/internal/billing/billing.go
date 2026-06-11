// Package billing is the public facade of the billing module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
//
// 設計 (ADR 0019): billing は Stripe Webhook 受信を起点として
// `payment` と `purchase` 両 context を呼び分けるオーケストレーター層。
// 依存方向: `billing → {payment, purchase}` (billing は両方を import する)。
//
// 命名は業務概念 (billing) で vendor 名 (stripe) ではない。
// 将来 PAY.JP 等を追加する場合も同一モジュール内に handler ファイルが並ぶだけで、
// 業務責務はここに閉じる (ADR 0019 ②)。
package billing

import (
	"net/http"

	billinghandler "github.com/Watari995/musclead/internal/billing/internal/handler"
	paymentpublicfunctions "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
)

// Module は billing module の公開 API。 現状は Handler のみ。
// 他 module から呼ばれる publicfunctions は今のところ無し。
type Module struct {
	Handler http.Handler
}

// NewModule は billing module を初期化する。 Composition Root (cmd/server/main.go) から呼ぶ。
//
// 想定する main.go 配線順 (ADR 0019 ④):
//
//	paymentModule  = payment.NewModule(...)
//	purchaseModule = purchase.NewModule(paymentModule.Command(), ...)
//	billingModule  = billing.NewModule(
//	    paymentModule.WebhookCommand(),
//	    paymentModule.StripeProcessor(),
//	    purchaseModule.Command(),
//	)
func NewModule(
	paymentWebhookCommand paymentpublicfunctions.PaymentWebhookCommand,
	stripeProcessor paymentpublicfunctions.StripeWebhookProcessor,
	purchaseCommand purchasepublicfunctions.PurchaseCommand,
) *Module {
	handler := billinghandler.NewWebhookHandler(paymentWebhookCommand, stripeProcessor, purchaseCommand)
	return &Module{
		Handler: handler,
	}
}
