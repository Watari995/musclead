package purchaseinfra

import (
	"context"
	"database/sql"
	"errors"

	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type subscriptionRepository struct {
	dbmap *gorp.DbMap
}

func NewSubscriptionRepository(dbmap *gorp.DbMap) purchasedomain.SubscriptionRepository {
	return &subscriptionRepository{dbmap: dbmap}
}

func (r *subscriptionRepository) FindLatestByUserID(ctx context.Context, userID valueobject.UserID) (*purchasedomain.Subscription, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row SubscriptionModel
	err = q.SelectOne(&row, "SELECT id, user_id, plan, status, subscription_order_id, payment_id, activated_at, expires_at, canceled_at, created_at, updated_at FROM subscriptions WHERE user_id = ? ORDER BY created_at DESC LIMIT 1", bytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toSubscription(row)
}

func (r *subscriptionRepository) FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*purchasedomain.Subscription, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := paymentID.Bytes()
	if err != nil {
		return nil, err
	}
	var row SubscriptionModel
	err = q.SelectOne(&row, "SELECT id, user_id, plan, status, subscription_order_id, payment_id, activated_at, expires_at, canceled_at, created_at, updated_at FROM subscriptions WHERE payment_id = ?", bytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toSubscription(row)
}

const upsertSubscriptionSQL = `
INSERT INTO subscriptions (id, user_id, plan, status, subscription_order_id, payment_id, activated_at, expires_at, canceled_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    status = VALUES(status),
    expires_at = VALUES(expires_at),
    canceled_at = VALUES(canceled_at),
    updated_at = VALUES(updated_at)
`

func (r *subscriptionRepository) Save(ctx context.Context, subscription *purchasedomain.Subscription) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertSubscriptionParams(subscription)
	if err != nil {
		return err
	}
	_, err = q.Exec(upsertSubscriptionSQL, params...)
	if err != nil {
		return err
	}
	return nil
}

func buildUpsertSubscriptionParams(subscription *purchasedomain.Subscription) ([]any, error) {
	bytes, err := subscription.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := subscription.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	plan := subscription.Plan().Value()
	status := subscription.Status().Value()
	subscriptionOrderIDBytes, err := sqlconv.NewBytesFromNullablePrimaryID(subscription.SubscriptionOrderID())
	if err != nil {
		return nil, err
	}
	paymentIDBytes, err := subscription.PaymentID().Bytes()
	if err != nil {
		return nil, err
	}
	activatedAt := subscription.ActivatedAt()
	expiresAt := subscription.ExpiresAt()
	canceledAt := sqlconv.ToNullTime(subscription.CanceledAt())
	createdAt := subscription.CreatedAt()
	updatedAt := subscription.UpdatedAt()
	return []any{bytes, userIDBytes, plan, status, subscriptionOrderIDBytes, paymentIDBytes, activatedAt, expiresAt, canceledAt, createdAt, updatedAt}, nil
}

func toSubscription(row SubscriptionModel) (*purchasedomain.Subscription, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.SubscriptionID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	plan, err := valueobject.NewSubscriptionPlanFromString(row.Plan)
	if err != nil {
		return nil, err
	}
	status, err := valueobject.NewSubscriptionStatusFromString(row.Status)
	if err != nil {
		return nil, err
	}
	subscriptionOrderID, err := sqlconv.NewPrimaryIDFromNullableBytes[valueobject.SubscriptionOrderID](row.SubscriptionOrderID)
	if err != nil {
		return nil, err
	}
	paymentID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.PaymentID](row.PaymentID)
	if err != nil {
		return nil, err
	}
	activatedAt := row.ActivatedAt
	expiresAt := row.ExpiresAt
	canceledAt := sqlconv.FromNullTime(row.CanceledAt)
	return purchasedomain.NewSubscription(*id, *userID, *plan, *status, subscriptionOrderID, *paymentID, activatedAt, expiresAt, canceledAt, row.CreatedAt, row.UpdatedAt), nil
}
