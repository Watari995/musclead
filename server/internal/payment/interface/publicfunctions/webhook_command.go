// Package publicfunctions の Webhook 系コマンド定義。
//
// 設計 (ADR 0019): billing module が Stripe Webhook を受信した時に、
// payment 側の状態遷移を起こすための公開 API。 billing handler から呼ばれる。
//
// `PaymentCommand` (InitiatePayment) は purchase 側からの「申込開始」 用、
// 本 interface はそれと責務が異なるため別 interface として分離する。
package publicfunctions

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

// CompletePaymentRequest は Stripe 'checkout.session.completed' 受信時の入力。
type CompletePaymentRequest struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata
}

// CompletePaymentResponse は billing handler が後段の
// `purchase.PurchaseCommand.ActivateSubscription` を呼ぶために必要な情報。
type CompletePaymentResponse struct {
	PaymentID valueobject.PaymentID
	UserID    valueobject.UserID
	Plan      valueobject.SubscriptionPlan
	ExpiresAt valueobject.Metadata // current_period_end など、 確定後に埋める
}

// CancelPaymentRequest は Stripe 'customer.subscription.deleted' 受信時の入力。
type CancelPaymentRequest struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata
}

// RenewPaymentRequest は Stripe 'invoice.paid' (月次更新) 受信時の入力。
type RenewPaymentRequest struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata
}

// HandleFailureRequest は処理不能 event を RetryStrategy に委譲するための入力。
type HandleFailureRequest struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata
}

// PaymentWebhookCommand は billing module (Webhook orchestrator) 専用の公開 API。
//
// 設計メモ (ADR 0018, 0019):
//   - 各メソッドは内部で独自 TX を張る (handler は TX を持たない)
//   - 冪等性は stripe_events UNIQUE で吸収、 重複受信は no-op
//   - CompletePayment のみ response を返す (billing が purchase.ActivateSubscription に流すため)
type PaymentWebhookCommand interface {
	CompletePayment(ctx context.Context, req CompletePaymentRequest) (CompletePaymentResponse, error)
	CancelPayment(ctx context.Context, req CancelPaymentRequest) error
	RenewPayment(ctx context.Context, req RenewPaymentRequest) error
	HandleFailure(ctx context.Context, req HandleFailureRequest) error
}
