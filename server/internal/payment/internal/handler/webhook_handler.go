package paymenthandler

import (
	"net/http"

	paymentusecase "github.com/Watari995/musclead/internal/payment/internal/usecase"
)

// WebhookHandler は POST /payment/webhook を処理する。
//
// 設計 (ADR 0018):
//   - handler は HTTP 入出力 + event_type による usecase dispatch のみ (薄い)
//   - 各 usecase が TX 内で stripe_events + 本処理を atomic に実行
//   - 失敗時は HandleFailure usecase 経由で RetryStrategy.OnFailure に委譲
//
// 流れ:
//  1. body / Stripe-Signature header 取得
//  2. parseWebhookEvent.Execute (TX 外、 検証 + パース)
//  3. event_type で分岐 → completePayment / cancelPayment / renewPayment
//  4. 失敗時は handleFailure.Execute → 5xx (Stripe 自動リトライ)
//  5. 成功は 200
type WebhookHandler struct {
	parseWebhookEvent *paymentusecase.ParseWebhookEvent
	completePayment   *paymentusecase.CompletePayment
	cancelPayment     *paymentusecase.CancelPayment
	renewPayment      *paymentusecase.RenewPayment
	handleFailure     *paymentusecase.HandleFailure
}

func NewWebhookHandler(
	parseWebhookEvent *paymentusecase.ParseWebhookEvent,
	completePayment *paymentusecase.CompletePayment,
	cancelPayment *paymentusecase.CancelPayment,
	renewPayment *paymentusecase.RenewPayment,
	handleFailure *paymentusecase.HandleFailure,
) *WebhookHandler {
	return &WebhookHandler{
		parseWebhookEvent: parseWebhookEvent,
		completePayment:   completePayment,
		cancelPayment:     cancelPayment,
		renewPayment:      renewPayment,
		handleFailure:     handleFailure,
	}
}

// Handle は POST /payment/webhook を処理する。
//
// TODO (User 実装):
//
//	body, _ := io.ReadAll(r.Body)
//	sig := r.Header.Get("Stripe-Signature")
//
//	event, err := h.parseWebhookEvent.Execute(ctx, paymentdomain.ParseWebhookEventInput{
//	    Payload: body, SignatureHeader: sig,
//	})
//	if err != nil { http.Error(w, "invalid signature", http.StatusUnauthorized); return }
//
//	var processErr error
//	switch event.EventType {
//	case "checkout.session.completed":
//	    processErr = h.completePayment.Execute(ctx, paymentusecase.CompletePaymentInput{...})
//	case "customer.subscription.deleted":
//	    processErr = h.cancelPayment.Execute(ctx, paymentusecase.CancelPaymentInput{...})
//	case "invoice.payment_succeeded":
//	    processErr = h.renewPayment.Execute(ctx, paymentusecase.RenewPaymentInput{...})
//	}
//
//	if processErr != nil {
//	    if err := h.handleFailure.Execute(ctx, paymentusecase.HandleFailureInput{Cause: processErr}); err != nil {
//	        http.Error(w, "internal error", http.StatusInternalServerError); return
//	    }
//	}
//	w.WriteHeader(http.StatusOK)
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "webhook handler not implemented", http.StatusNotImplemented)
}
