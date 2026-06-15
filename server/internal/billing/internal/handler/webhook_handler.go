package billinghandler

import (
	"io"
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	paymentpublicfunctions "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/shared/httpx"
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
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	// Stripe-Signature を取得
	signatureHeader := r.Header.Get("Stripe-Signature")
	if signatureHeader == "" {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("Stripe-Signature is required"))
		return
	}
	event, err := h.stripeProcessor.ParseAndVerify(r.Context(), payload, signatureHeader)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	// event.EventTypeで分岐
	var errResponse error
	switch event.EventType {
	case "checkout.session.completed":
		var resp paymentpublicfunctions.CompletePaymentResponse
		resp, errResponse = h.paymentWebhookCommand.CompletePayment(r.Context(), paymentpublicfunctions.CompletePaymentRequest{
			StripeEventID: event.StripeEventID,
			EventType:     event.EventType,
			Payload:       event.Payload,
		})
		// completePaymentが成功した時のみActivateSubscriptionを呼ぶ
		if errResponse == nil {
			errResponse = h.purchaseCommand.ActivateSubscription(r.Context(), purchasepublicfunctions.ActivateSubscriptionRequest{
				PaymentID: resp.PaymentID,
				UserID:    resp.UserID,
				Plan:      resp.Plan,
				ExpiresAt: resp.ExpiresAt,
			})
		}
	case "customer.subscription.deleted":
		var resp paymentpublicfunctions.CancelPaymentResponse
		resp, errResponse = h.paymentWebhookCommand.CancelPayment(r.Context(), paymentpublicfunctions.CancelPaymentRequest{
			StripeEventID: event.StripeEventID,
			EventType:     event.EventType,
			Payload:       event.Payload,
		})
		// cancelPaymentが成功した時のみCancelSubscriptionを呼ぶ
		if errResponse == nil {
			errResponse = h.purchaseCommand.CancelSubscription(r.Context(), purchasepublicfunctions.CancelSubscriptionRequest{
				PaymentID: resp.PaymentID,
			})
		}
	case "invoice.paid":
		var resp paymentpublicfunctions.RenewPaymentResponse
		resp, errResponse = h.paymentWebhookCommand.RenewPayment(r.Context(), paymentpublicfunctions.RenewPaymentRequest{
			StripeEventID: event.StripeEventID,
			EventType:     event.EventType,
			Payload:       event.Payload,
		})
		// renewPaymentが成功した時のみrenewSubscriptionを呼ぶ
		if errResponse == nil {
			errResponse = h.purchaseCommand.RenewSubscription(r.Context(), purchasepublicfunctions.RenewSubscriptionRequest{
				PaymentID: resp.PaymentID,
				ExpiresAt: resp.ExpiresAt,
			})
		}
	default:
		httpx.WriteOK(w)
		return
	}
	// 全体のswitchでどこかで失敗したらHandleFailureを呼ぶ
	if errResponse != nil {
		if ferr := h.paymentWebhookCommand.HandleFailure(r.Context(), paymentpublicfunctions.HandleFailureRequest{
			StripeEventID: event.StripeEventID,
			EventType:     event.EventType,
			Payload:       event.Payload,
			Cause:         errResponse,
		}); ferr != nil {
			httpx.WriteError(w, ferr)
			return
		}
	}
	// 成功で200レスポンスを返す
	httpx.WriteOK(w)
}
