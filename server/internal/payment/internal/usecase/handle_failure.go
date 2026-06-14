package paymentusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
)

// HandleFailure は Webhook 処理失敗時の挙動を RetryStrategy に委譲する薄い wrapper usecase。
//
// 設計 (ADR 0014, 0018):
//   - handler が domain interface (RetryStrategy) を直接呼ばないようにする境界
//   - MVP は ExternalRetryStrategy = err を return → handler が 500 を返す → Stripe 自動リトライ
//   - 将来 SelfManagedRetryStrategy = failed_webhook_events に記録 + nil → 200 → 自前リトライ
type HandleFailure struct {
	retryStrategy paymentdomain.RetryStrategy
}

func (uc *HandleFailure) HandleFailure(ctx context.Context, input publicfunctions.HandleFailureRequest) error {
	stripeEvent := paymentdomain.CreateStripeEvent(input.StripeEventID, input.EventType, input.Payload)
	return uc.retryStrategy.OnFailure(ctx, stripeEvent, input.Cause)
}

func NewHandleFailure(retryStrategy paymentdomain.RetryStrategy) *HandleFailure {
	return &HandleFailure{retryStrategy: retryStrategy}
}
