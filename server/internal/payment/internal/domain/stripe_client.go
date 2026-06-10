package paymentdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateCustomerInput struct {
	UserID valueobject.UserID
	Email  valueobject.Email
}

type CreateCheckoutSessionInput struct {
	CustomerID string
	PriceID    string
	PaymentID  valueobject.PaymentID
}

type CreateCheckoutSessionOutput struct {
	SessionID          string
	CheckoutSessionURL valueobject.URL
}

// ParseWebhookEventInput は Stripe Webhook の生 payload と署名ヘッダを受け取る。
type ParseWebhookEventInput struct {
	Payload         []byte // HTTP body 全文 (改ざん検証のため []byte で保持、 文字列化しない)
	SignatureHeader string // Stripe-Signature ヘッダの値
}

// ParseWebhookEventOutput は検証済みの Stripe event 情報。
// EventType に応じて usecase で分岐し、 Payload から必要な field を取り出す。
type ParseWebhookEventOutput struct {
	StripeEventID string               // Stripe 発行の一意 ID (evt_xxx)、 stripe_events.stripe_event_id に保存
	EventType     string               // 'checkout.session.completed' 等、 usecase 分岐に使う
	Payload       valueobject.Metadata // event の data 部分、 stripe_events.payload に保存
}

// StripeClient は Stripe API への呼び出しを抽象化する。
// musclead の usecase / domain は Stripe SDK を直接 import しない (ACL)。
type StripeClient interface {
	CreateCustomer(ctx context.Context, input CreateCustomerInput) (customerID string, err error)
	CreateCheckoutSession(ctx context.Context, input CreateCheckoutSessionInput) (output CreateCheckoutSessionOutput, err error)

	// ParseWebhookEvent は Stripe Webhook の署名検証 + event 取り出しを 1 度に行う。
	// 署名不一致 / パース失敗時は error を返す。
	// 検証は TX 外で実行する想定 (DB 操作なし、 純粋な HMAC + JSON 処理)。
	ParseWebhookEvent(ctx context.Context, input ParseWebhookEventInput) (output ParseWebhookEventOutput, err error)
}
