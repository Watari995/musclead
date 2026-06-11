package billinghandler

import (
	"net/http"

	paymentpublicfunctions "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
)

// WebhookHandler は Stripe Webhook を受信し、 payment / purchase 両 context に dispatch する。
//
// 設計 (ADR 0014 / 0018 / 0019):
//   - handler 自身が「両 context のオーケストレーター」 (ADR 0014 末尾)
//   - 単一 TX は張らない、 各 usecase が独自 TX (ADR 0019 ⑥)
//   - 冪等性は stripe_events UNIQUE + ActivateSubscription FindByPaymentID 既存チェックで吸収
type WebhookHandler struct {
	paymentWebhookCommand paymentpublicfunctions.PaymentWebhookCommand
	stripeProcessor       paymentpublicfunctions.StripeWebhookProcessor
	purchaseCommand       purchasepublicfunctions.PurchaseCommand
}

// Handle godoc
//
// @Summary  Stripe Webhook 受信 (Pro 申込み完了 / 解約 / 月次更新)
// @Tags     billing
// @Accept   json
// @Produce  json
// @Param    Stripe-Signature header string true "Stripe 署名"
// @Success  200
// @Failure  400 {object} httpx.ErrorResponse
// @Router   /billing/webhook [post]
//
// 流れ (ADR 0019):
//  1. r.Body を読み取り、 r.Header.Get("Stripe-Signature") を取得
//  2. h.stripeProcessor.ParseAndVerify で署名検証 + StripeEvent 化 (TX 外)
//  3. event.EventType で dispatch:
//     - 'checkout.session.completed'   → CompletePayment → ActivateSubscription
//     - 'customer.subscription.deleted' → CancelPayment (purchase 側 Cancel は将来)
//     - 'invoice.paid'                  → RenewPayment   (purchase 側 Renew は将来)
//     - その他                           → HandleFailure
//  4. 各段階の error は httpx.WriteError で HTTP status に変換、 5xx は Stripe リトライ任せ (ADR 0014 ⑤)
//
// 設計メモ:
//   - publicfunctions のみ import (`payment/internal/*`, `purchase/internal/*` は触らない)
//   - CompletePayment レスポンスから ActivateSubscriptionRequest を組み立てて purchase に渡す
//   - dispatch ロジックは Stripe プロトコル解釈なので handler 配置で OK (ADR 0018 ①)
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
}

// NewWebhookHandler は WebhookHandler を組み立て、 POST /billing/webhook にマウントした http.Handler を返す。
// 各 module の publicfunctions のみ受け取り、 internal の型は知らない (ADR 0019 ①)。
func NewWebhookHandler(
	paymentWebhookCommand paymentpublicfunctions.PaymentWebhookCommand,
	stripeProcessor paymentpublicfunctions.StripeWebhookProcessor,
	purchaseCommand purchasepublicfunctions.PurchaseCommand,
) http.Handler {
	h := &WebhookHandler{
		paymentWebhookCommand: paymentWebhookCommand,
		stripeProcessor:       stripeProcessor,
		purchaseCommand:       purchaseCommand,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /billing/webhook", h.Handle)
	return mux
}
