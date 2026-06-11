package purchasedomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

// SubscriptionRepository は subscriptions テーブルへの永続化を抽象化する。
//
// 設計:
//   - 「Not Found なら (nil, nil)」 musclead 既存流儀
//   - Save は INSERT / UPDATE 兼用 (upsert)
//
// TODO (User 実装):
//   - メソッド候補:
//     Save(ctx, sub *Subscription) error
//     FindLatestByUserID(ctx, userID valueobject.UserID) (*Subscription, error)
//     → Pro 判定 / マイページ表示用、 user 1 人につき最新 1 件を返す
//     FindByPaymentID(ctx, paymentID valueobject.PaymentID) (*Subscription, error)
//     → Webhook で payment_id 経由で subscription を引きたい時
type SubscriptionRepository interface {
	FindLatestByUserID(ctx context.Context, userID valueobject.UserID) (*Subscription, error)
	FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*Subscription, error)
	Save(ctx context.Context, subscription *Subscription) error
}
