// Package publicfunctions は payment module が他 module に公開する Command / Query interface を定義する。
// 他 module (purchase 等) はこのパッケージにのみ依存する。
//
// 設計 (ADR 0013): Facade パターン、 SODA 流儀の「決済開始 / 完了」 公開。
package publicfunctions

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

// InitiatePaymentRequest は Pro 申込時に purchase 集約から渡される入力。
// 金額は Stripe 側の Price object で管理しているためここでは受け取らない。
type InitiatePaymentRequest struct {
	UserID  valueobject.UserID
	Email   valueobject.Email
	PriceID string
}

// InitiatePaymentResponse は申込開始の結果 (Stripe Checkout への遷移情報)。
type InitiatePaymentResponse struct {
	PaymentID   valueobject.PaymentID
	CheckoutURL valueobject.URL
}

// PaymentCommand は payment 集約に対する書き込み系操作の公開 API。
// purchase context からは Module.Command() 経由でアクセスする。
//
// MVP では InitiatePayment のみ公開。 CompletePayment / CancelPayment / RenewPayment は
// Webhook 受信時に内部の handler から呼ばれるため、 外部公開しない。
type PaymentCommand interface {
	InitiatePayment(ctx context.Context, req InitiatePaymentRequest) (InitiatePaymentResponse, error)
}
