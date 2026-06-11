package purchaseinfra

import (
	"context"
	"errors"

	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

// errSubscriptionOrderRepoNotImplemented は skeleton 状態のメソッドが返す sentinel error。
// User が中身を実装したら各メソッドから削除する。
var errSubscriptionOrderRepoNotImplemented = errors.New("subscription_order_repository: method not implemented")

type subscriptionOrderRepository struct {
	dbmap *gorp.DbMap
}

func NewSubscriptionOrderRepository(dbmap *gorp.DbMap) purchasedomain.SubscriptionOrderRepository {
	return &subscriptionOrderRepository{dbmap: dbmap}
}

func (r *subscriptionOrderRepository) FindPendingByUserID(ctx context.Context, userID valueobject.UserID) (*purchasedomain.SubscriptionOrder, error) {
	return nil, errSubscriptionOrderRepoNotImplemented
}

func (r *subscriptionOrderRepository) FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*purchasedomain.SubscriptionOrder, error) {
	return nil, errSubscriptionOrderRepoNotImplemented
}

func (r *subscriptionOrderRepository) Save(ctx context.Context, order *purchasedomain.SubscriptionOrder) error {
	return errSubscriptionOrderRepoNotImplemented
}

// 共通 SELECT 句。 field 順は migration (000020_create_subscription_orders.up.sql) の column 順に揃える。
//
// TODO (User 実装): const subscriptionOrderSelectColumns = ``
//   - id, user_id, plan, status, payment_id, succeeded_at, failed_at, created_at, updated_at

// TODO (User 実装): FindPendingByUserID
//   - 入力: user_id (BINARY(16) に変換、 user_id.Bytes())
//   - SQL: SELECT ... FROM subscription_orders WHERE user_id = ? AND status = 'pending' ORDER BY created_at DESC LIMIT 1
//   - 結果: sql.ErrNoRows なら (nil, nil) musclead 流儀
//   - 変換: toSubscriptionOrder(row)
// func (r *subscriptionOrderRepository) FindPendingByUserID(ctx context.Context, userID valueobject.UserID) (*purchasedomain.SubscriptionOrder, error)

// TODO (User 実装): FindByPaymentID
//   - 入力: payment_id (BINARY(16) に変換、 paymentID.Bytes())
//   - SQL: SELECT ... FROM subscription_orders WHERE payment_id = ?

// TODO (User 実装): Save (upsert)
//   - INSERT ... ON DUPLICATE KEY UPDATE
//   - status / payment_id / succeeded_at / failed_at / updated_at を UPDATE 句に
//   - 参考: payment_repository.go の upsertPaymentSQL

// TODO (User 実装): toSubscriptionOrder(row) (*SubscriptionOrder, error)
//   - id: sqlconv.NewPrimaryIDFromBytes[valueobject.SubscriptionOrderID](row.ID)
//   - paymentID: row.PaymentID が nil なら nil、 そうでなければ sqlconv.NewPrimaryIDFromBytes[valueobject.PaymentID]
//   - status: valueobject.NewSubscriptionOrderStatusFromString(row.Status) → *SubscriptionOrderStatus
//   - plan: valueobject.NewSubscriptionPlanFromString(row.Plan)
//   - succeededAt / failedAt: sqlconv.FromNullTime
//   - 全部揃ったら purchasedomain.NewSubscriptionOrder(...) で復元

