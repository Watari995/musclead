package paymentdomain

import (
	"context"
	"errors"
)

var ErrStripeEventAlreadyExists = errors.New("stripe event already exists")

// StripeEventRepository は stripe_events テーブルへの永続化を抽象化する。
//
// 設計 (ADR 0014):
//   - 受信した全 Stripe event を生 payload で残す (監査 / デバッグ / 再処理)
//   - stripe_event_id (Stripe 側の evt_xxx) は UNIQUE 制約で二重処理を物理的に防ぐ
type StripeEventRepository interface {
	// 同一 event の重複受信時に「既に処理したか」 を確認するために使う。
	FindByStripeEventID(ctx context.Context, stripeEventID string) (*StripeEvent, error)
	Create(ctx context.Context, event *StripeEvent) error
	UpdateProcessStatus(ctx context.Context, event *StripeEvent) error
}
