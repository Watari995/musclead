package purchasedomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

// SubscriptionOrderRepository は subscription_orders テーブルへの永続化を抽象化する。
//
// 設計:
//   - 「Not Found なら (nil, nil)」 musclead 既存流儀
//   - Save は INSERT / UPDATE 兼用 (upsert)
//
// TODO (User 実装):
//   - メソッド候補:
//     Save(ctx, order *SubscriptionOrder) error
//     FindByID(ctx, id valueobject.SubscriptionOrderID) (*SubscriptionOrder, error)
//     FindPendingByUserID(ctx, userID valueobject.UserID) (*SubscriptionOrder, error)
//     → 申込みリトライ時に既存 pending を再利用するか判定
//     FindByPaymentID(ctx, paymentID valueobject.PaymentID) (*SubscriptionOrder, error)
//     → Webhook で payment_id 経由で order を引きたい時
type SubscriptionOrderRepository interface {
	FindPendingByUserID(ctx context.Context, userID valueobject.UserID) (*SubscriptionOrder, error)
	FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*SubscriptionOrder, error)
	Save(ctx context.Context, order *SubscriptionOrder) error
}
