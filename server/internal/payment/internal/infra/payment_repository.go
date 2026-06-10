package paymentinfra

import (
	"context"
	"database/sql"
	"errors"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type paymentRepository struct {
	dbmap *gorp.DbMap
}

func NewPaymentRepository(dbmap *gorp.DbMap) paymentdomain.PaymentRepository {
	return &paymentRepository{dbmap: dbmap}
}

// 共通 SELECT 句。 field 順は migration (000015_create_payments.up.sql) の column 順に揃える。
const paymentSelectColumns = `
	id,
	user_id,
	amount,
	currency,
	status,
	stripe_customer_id,
	stripe_subscription_id,
	stripe_checkout_session_id,
	checkout_url,
	current_period_end,
	succeeded_at,
	failed_at,
	failure_reason,
	created_at,
	updated_at`

func (r *paymentRepository) FindByID(ctx context.Context, id valueobject.PaymentID) (*paymentdomain.Payment, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	var row PaymentModel
	err = q.SelectOne(&row,
		`SELECT`+paymentSelectColumns+` FROM payments WHERE id = ?`, bytes)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toPayment(row)
}

func (r *paymentRepository) FindLatestSucceededByUserID(ctx context.Context, userID valueobject.UserID) (*paymentdomain.Payment, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row PaymentModel
	err = q.SelectOne(&row,
		`SELECT`+paymentSelectColumns+` FROM payments WHERE user_id = ? AND status = ? ORDER BY created_at DESC LIMIT 1`,
		bytes, string(valueobject.PaymentStatusSucceeded))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toPayment(row)
}

func (r *paymentRepository) FindByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*paymentdomain.Payment, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var row PaymentModel
	err := q.SelectOne(&row,
		`SELECT`+paymentSelectColumns+` FROM payments WHERE stripe_subscription_id = ?`, stripeSubscriptionID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toPayment(row)
}

const upsertPaymentSQL = `
INSERT INTO payments (id, user_id, amount, currency, status, stripe_customer_id, stripe_subscription_id, stripe_checkout_session_id, checkout_url, current_period_end, succeeded_at, failed_at, failure_reason, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    amount = VALUES(amount),
    currency = VALUES(currency),
    status = VALUES(status),
    stripe_customer_id = VALUES(stripe_customer_id),
    stripe_subscription_id = VALUES(stripe_subscription_id),
    stripe_checkout_session_id = VALUES(stripe_checkout_session_id),
    checkout_url = VALUES(checkout_url),
    current_period_end = VALUES(current_period_end),
    succeeded_at = VALUES(succeeded_at),
    failed_at = VALUES(failed_at),
    failure_reason = VALUES(failure_reason),
    updated_at = VALUES(updated_at)
`

func (r *paymentRepository) Save(ctx context.Context, payment *paymentdomain.Payment) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertPaymentParams(payment)
	if err != nil {
		return err
	}
	_, err = q.Exec(upsertPaymentSQL, params...)
	if err != nil {
		return err
	}

	return nil
}

func toPayment(row PaymentModel) (*paymentdomain.Payment, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.PaymentID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	amount, err := valueobject.NewNonNegativeInt(row.Amount)
	if err != nil {
		return nil, err
	}
	currency, err := valueobject.NewCurrencyFromString(row.Currency)
	if err != nil {
		return nil, err
	}
	status, err := valueobject.NewPaymentStatusFromString(row.Status)
	if err != nil {
		return nil, err
	}
	checkoutURL, err := sqlconv.NewURLFromNullString(row.CheckoutURL)
	if err != nil {
		return nil, err
	}
	currentPeriodEnd := sqlconv.FromNullTime(row.CurrentPeriodEnd)
	return paymentdomain.NewPayment(
		*id,
		*userID,
		*amount,
		*currency,
		*status,
		sqlconv.NewStringFromNullString(row.StripeCustomerID),
		sqlconv.NewStringFromNullString(row.StripeSubscriptionID),
		sqlconv.NewStringFromNullString(row.StripeCheckoutSessionID),
		checkoutURL,
		currentPeriodEnd,
		sqlconv.FromNullTime(row.SucceededAt),
		sqlconv.FromNullTime(row.FailedAt),
		sqlconv.NewStringFromNullString(row.FailureReason),
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func buildUpsertPaymentParams(payment *paymentdomain.Payment) ([]any, error) {
	bytes, err := payment.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := payment.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	amount := payment.Amount().Value()
	currency := payment.Currency().Value()
	status := payment.Status().Value()
	stripeCustomerID := sqlconv.StringPtrToNullString(payment.StripeCustomerID())
	stripeSubscriptionID := sqlconv.StringPtrToNullString(payment.StripeSubscriptionID())
	stripeCheckoutSessionID := sqlconv.StringPtrToNullString(payment.StripeCheckoutSessionID())
	checkoutURL := sqlconv.URLPtrToNullString(payment.CheckoutURL())
	currentPeriodEnd := sqlconv.ToNullTime(payment.CurrentPeriodEnd())
	succeededAt := sqlconv.ToNullTime(payment.SucceededAt())
	failedAt := sqlconv.ToNullTime(payment.FailedAt())
	failureReason := sqlconv.StringPtrToNullString(payment.FailureReason())
	createdAt := payment.CreatedAt()
	updatedAt := payment.UpdatedAt()
	return []any{
		bytes,
		userIDBytes,
		amount,
		currency,
		status,
		stripeCustomerID,
		stripeSubscriptionID,
		stripeCheckoutSessionID,
		checkoutURL,
		currentPeriodEnd,
		succeededAt,
		failedAt,
		failureReason,
		createdAt,
		updatedAt,
	}, nil
}
