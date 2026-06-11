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

type subscriptionOrderRepository struct {
	dbmap *gorp.DbMap
}

func NewSubscriptionOrderRepository(dbmap *gorp.DbMap) purchasedomain.SubscriptionOrderRepository {
	return &subscriptionOrderRepository{dbmap: dbmap}
}

func (r *subscriptionOrderRepository) FindPendingByUserID(ctx context.Context, userID valueobject.UserID) (*purchasedomain.SubscriptionOrder, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row SubscriptionOrderModel
	err = q.SelectOne(&row, "SELECT id, user_id, plan, status, payment_id, succeeded_at, failed_at, created_at, updated_at FROM subscription_orders WHERE user_id = ? AND status = ? LIMIT 1", bytes, string(valueobject.SubscriptionOrderStatusPending))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toSubscriptionOrder(row)
}

func (r *subscriptionOrderRepository) FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*purchasedomain.SubscriptionOrder, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := paymentID.Bytes()
	if err != nil {
		return nil, err
	}
	var row SubscriptionOrderModel
	err = q.SelectOne(&row, "SELECT id, user_id, plan, status, payment_id, succeeded_at, failed_at, created_at, updated_at FROM subscription_orders WHERE payment_id = ? LIMIT 1", bytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toSubscriptionOrder(row)
}

const upsertSubscriptionOrderSQL = `
INSERT INTO subscription_orders (id, user_id, plan, status, payment_id, succeeded_at, failed_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    user_id = VALUES(user_id),
    plan = VALUES(plan),
    status = VALUES(status),
    payment_id = VALUES(payment_id),
    succeeded_at = VALUES(succeeded_at),
    failed_at = VALUES(failed_at),
    updated_at = VALUES(updated_at)
`

func (r *subscriptionOrderRepository) Save(ctx context.Context, order *purchasedomain.SubscriptionOrder) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertSubscriptionOrderParams(order)
	if err != nil {
		return err
	}
	_, err = q.Exec(upsertSubscriptionOrderSQL, params...)
	if err != nil {
		return err
	}
	return nil
}

func buildUpsertSubscriptionOrderParams(order *purchasedomain.SubscriptionOrder) ([]any, error) {
	bytes, err := order.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := order.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	plan := order.Plan().Value()
	status := order.Status().Value()
	paymentIDBytes, err := sqlconv.NewBytesFromNullablePrimaryID(order.PaymentID())
	if err != nil {
		return nil, err
	}
	succeededAt := sqlconv.ToNullTime(order.SucceededAt())
	failedAt := sqlconv.ToNullTime(order.FailedAt())
	createdAt := order.CreatedAt()
	updatedAt := order.UpdatedAt()
	return []any{bytes, userIDBytes, plan, status, paymentIDBytes, succeededAt, failedAt, createdAt, updatedAt}, nil
}

func toSubscriptionOrder(row SubscriptionOrderModel) (*purchasedomain.SubscriptionOrder, error) {
	// nullStringなどに注意して変換する
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.SubscriptionOrderID](row.ID)
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
	status, err := valueobject.NewSubscriptionOrderStatusFromString(row.Status)
	if err != nil {
		return nil, err
	}
	paymentID, err := sqlconv.NewPrimaryIDFromNullableBytes[valueobject.PaymentID](row.PaymentID)
	if err != nil {
		return nil, err
	}
	succeededAt := sqlconv.FromNullTime(row.SucceededAt)
	failedAt := sqlconv.FromNullTime(row.FailedAt)

	return purchasedomain.NewSubscriptionOrder(*id, *userID, *plan, *status, paymentID, succeededAt, failedAt, row.CreatedAt, row.UpdatedAt), nil
}
