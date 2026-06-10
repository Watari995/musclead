package paymentusecase

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
)

// ParseWebhookEvent は Stripe Webhook の署名検証 + event 取り出しを行う薄い wrapper usecase。
//
// 設計 (ADR 0018):
//   - handler が domain interface (StripeClient) を直接呼ばないようにする境界
//   - 中身は stripeClient.ParseWebhookEvent を呼ぶだけ
//   - TX 外で実行 (DB 触らない、 純粋な HMAC + JSON 処理)
type ParseWebhookEvent struct {
	stripeClient paymentdomain.StripeClient
}

func NewParseWebhookEvent(stripeClient paymentdomain.StripeClient) *ParseWebhookEvent {
	return &ParseWebhookEvent{stripeClient: stripeClient}
}

// Execute は payload + signature を受け取り、 検証済み event 情報を返す。
//
// TODO (User 実装):
//
//	return uc.stripeClient.ParseWebhookEvent(ctx, input)
func (uc *ParseWebhookEvent) Execute(ctx context.Context, input paymentdomain.ParseWebhookEventInput) (paymentdomain.ParseWebhookEventOutput, error) {
	return paymentdomain.ParseWebhookEventOutput{}, errNotImplemented
}
