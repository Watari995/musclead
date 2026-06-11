package paymenthandler

import (
	"io"
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	paymentusecase "github.com/Watari995/musclead/internal/payment/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/httpx"
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

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("failed to read request body").Wrap(err))
		return
	}
	signature := r.Header.Get("Stripe-Signature")

	// 1. parse + 署名検証
	event, err := h.parseWebhookEvent.Execute(r.Context(), paymentusecase.ParseWebhookEventInput{
		Payload:         body,
		SignatureHeader: signature,
	})
	if err != nil {
		httpx.WriteError(w, myerror.NewUnauthorizedError().SetMessage("invalid signature").Wrap(err))
		return
	}

	// 2. event_type で分岐 (Stripe Webhook プロトコルの解釈)
	var processErr error
	switch event.EventType {
	case "checkout.session.completed":
		processErr = h.completePayment.Execute(r.Context(), paymentusecase.CompletePaymentInput{
			StripeEventID: event.StripeEventID,
			EventType:     event.EventType,
			Payload:       event.Payload,
		})
	case "customer.subscription.deleted":
		processErr = h.cancelPayment.Execute(r.Context(), paymentusecase.CancelPaymentInput{
			StripeEventID: event.StripeEventID,
			EventType:     event.EventType,
			Payload:       event.Payload,
		})
	case "invoice.payment_succeeded":
		processErr = h.renewPayment.Execute(r.Context(), paymentusecase.RenewPaymentInput{
			StripeEventID: event.StripeEventID,
			EventType:     event.EventType,
			Payload:       event.Payload,
		})
	default:
		// musclead が興味ない event_type (customer.created 等) は無視して 200 を返す (Stripe のリトライ抑止)
	}

	// 3. 失敗時は RetryStrategy に委譲 (ExternalRetryStrategy なら err を返して 500 → Stripe 自動リトライ)
	if processErr != nil {
		if err := h.handleFailure.Execute(r.Context(), paymentusecase.HandleFailureInput{
			Cause: processErr,
		}); err != nil {
			httpx.WriteError(w, myerror.NewInternalError().Wrap(err))
			return
		}
	}

	httpx.WriteOK(w)
}
