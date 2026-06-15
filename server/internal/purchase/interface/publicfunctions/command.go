// Package publicfunctions は purchase module が他 module に公開する Command / Query interface を定義する。
//
// 設計 (ADR 0019): billing module (Webhook orchestrator) が Stripe イベント受信時に
// 「Pro 化を確定させる」 ために本 interface を呼ぶ。 依存方向は `billing → purchase`。
package publicfunctions

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// ActivateSubscriptionRequest は Stripe checkout 完了時に billing handler から渡される入力。
//
// 設計メモ:
//   - PaymentID で `subscriptions.payment_id` を引いて冪等性チェック (ADR 0014 ③)
//   - ExpiresAt は Stripe 側の current_period_end を渡す (purchase は Stripe 知らない)
type ActivateSubscriptionRequest struct {
	PaymentID valueobject.PaymentID
	UserID    valueobject.UserID
	Plan      valueobject.SubscriptionPlan
	ExpiresAt time.Time
}

type RenewSubscriptionRequest struct {
	PaymentID valueobject.PaymentID
	ExpiresAt time.Time
}

type CancelSubscriptionRequest struct {
	PaymentID valueobject.PaymentID
}

// PurchaseCommand は purchase 集約に対する書き込み系操作の公開 API。
//
// MVP (Phase 2 後半) では ActivateSubscription のみ公開。
// 将来 Webhook で「解約」 「期限切れ」 を扱う際は Cancel / Expire を追加。
type PurchaseCommand interface {
	ActivateSubscription(ctx context.Context, req ActivateSubscriptionRequest) error
	RenewSubscription(ctx context.Context, req RenewSubscriptionRequest) error
	CancelSubscription(ctx context.Context, req CancelSubscriptionRequest) error
}
