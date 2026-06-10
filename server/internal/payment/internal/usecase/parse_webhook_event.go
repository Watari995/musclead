package paymentusecase

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ParseWebhookEventInput struct {
	Payload         []byte // HTTP body 全文 (改ざん検証のため []byte で保持、 文字列化しない)
	SignatureHeader string // Stripe-Signature ヘッダの値
}

type ParseWebhookEventOutput struct {
	StripeEventID string               // Stripe 発行の一意 ID (evt_xxx)、 stripe_events.stripe_event_id に保存
	EventType     string               // 'checkout.session.completed' 等、 usecase 分岐に使う
	Payload       valueobject.Metadata // event の data 部分、 stripe_events.payload に保存
}

// ParseWebhookEvent は Stripe Webhook の署名検証 + event 取り出しを行う薄い wrapper usecase。
//
// 設計 (ADR 0018):
//   - handler が domain interface (StripeClient) を直接呼ばないようにする境界
//   - 中身は stripeClient.ParseWebhookEvent を呼ぶだけ
//   - TX 外で実行 (DB 触らない、 純粋な HMAC + JSON 処理)
type ParseWebhookEvent struct {
	stripeClient paymentdomain.StripeClient
}

func (uc *ParseWebhookEvent) Execute(ctx context.Context, input ParseWebhookEventInput) (ParseWebhookEventOutput, error) {
	output, err := uc.stripeClient.ParseWebhookEvent(ctx, paymentdomain.ParseWebhookEventInput{
		Payload:         input.Payload,
		SignatureHeader: input.SignatureHeader,
	})
	if err != nil {
		return ParseWebhookEventOutput{}, err
	}
	return ParseWebhookEventOutput{
		StripeEventID: output.StripeEventID,
		EventType:     output.EventType,
		Payload:       output.Payload,
	}, nil
}

func NewParseWebhookEvent(stripeClient paymentdomain.StripeClient) *ParseWebhookEvent {
	return &ParseWebhookEvent{stripeClient: stripeClient}
}
