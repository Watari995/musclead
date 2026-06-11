package purchasedomain

import (
	"context"
)

// SubscriptionRepository は subscriptions テーブルへの永続化を抽象化する。
//
// 設計:
//   - 「Not Found なら (nil, nil)」 musclead 既存流儀
//   - Save は INSERT / UPDATE 兼用 (upsert)
//
// TODO (User 実装):
//   - メソッド候補:
//       Save(ctx, sub *Subscription) error
//       FindLatestByUserID(ctx, userID valueobject.UserID) (*Subscription, error)
//           → Pro 判定 / マイページ表示用、 user 1 人につき最新 1 件を返す
//       FindByPaymentID(ctx, paymentID valueobject.PaymentID) (*Subscription, error)
//           → Webhook で payment_id 経由で subscription を引きたい時
type SubscriptionRepository interface {
	// Placeholder: User がメソッドを定義したら削除
	_placeholder(ctx context.Context) error
}
