package purchaseinfra

import (
	"context"
	"errors"

	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

// errSubscriptionRepoNotImplemented は skeleton 状態のメソッドが返す sentinel error。
// User が中身を実装したら各メソッドから削除する。
var errSubscriptionRepoNotImplemented = errors.New("subscription_repository: method not implemented")

type subscriptionRepository struct {
	dbmap *gorp.DbMap
}

func NewSubscriptionRepository(dbmap *gorp.DbMap) purchasedomain.SubscriptionRepository {
	return &subscriptionRepository{dbmap: dbmap}
}

func (r *subscriptionRepository) FindLatestByUserID(ctx context.Context, userID valueobject.UserID) (*purchasedomain.Subscription, error) {
	return nil, errSubscriptionRepoNotImplemented
}

func (r *subscriptionRepository) FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*purchasedomain.Subscription, error) {
	return nil, errSubscriptionRepoNotImplemented
}

func (r *subscriptionRepository) Save(ctx context.Context, subscription *purchasedomain.Subscription) error {
	return errSubscriptionRepoNotImplemented
}

// 共通 SELECT 句。 field 順は migration (000021_create_subscriptions.up.sql) の column 順に揃える。
//
// TODO (User 実装): const subscriptionSelectColumns = ``
//   - id, user_id, plan, status, subscription_order_id, payment_id,
//     activated_at, expires_at, canceled_at, created_at, updated_at

// TODO (User 実装): FindLatestByUserID
//   - 入力: user_id (BINARY(16))
//   - SQL: SELECT ... FROM subscriptions WHERE user_id = ? ORDER BY created_at DESC LIMIT 1
//   - 結果: sql.ErrNoRows なら (nil, nil)
//   - 用途: Pro 判定 (subscription.expires_at > NOW() で gate)、 マイページ表示

// TODO (User 実装): FindByPaymentID
//   - SQL: WHERE payment_id = ?
//   - 用途: Webhook (renew / cancel) で payment 経由 subscription を引く

// TODO (User 実装): Save (upsert)
//   - INSERT ... ON DUPLICATE KEY UPDATE
//   - UPDATE 句: status / expires_at / canceled_at / updated_at が主、 他は変えない想定
//   - 参考: payment_repository.go の upsertPaymentSQL

// TODO (User 実装): toSubscription(row) (*Subscription, error)
//   - id, userID, paymentID: sqlconv.NewPrimaryIDFromBytes[T]
//   - subscriptionOrderID: row.SubscriptionOrderID が nil なら nil
//   - plan: valueobject.NewSubscriptionPlanFromString
//   - status: valueobject.NewSubscriptionStatusFromString
//   - canceledAt: sqlconv.FromNullTime
//   - 全部揃ったら purchasedomain.NewSubscription(...) で復元

