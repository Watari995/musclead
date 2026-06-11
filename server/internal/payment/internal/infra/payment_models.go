package paymentinfra

import (
	"database/sql"
	"time"
)

// PaymentModel は payments テーブルの gorp マッピング。
// migration: sql/migrations/000015_create_payments.up.sql
type PaymentModel struct {
	ID                      []byte         `db:"id"`
	UserID                  []byte         `db:"user_id"`
	Currency                string         `db:"currency"`
	Status                  string         `db:"status"`
	StripeCustomerID        sql.NullString `db:"stripe_customer_id"`
	StripeSubscriptionID    sql.NullString `db:"stripe_subscription_id"`
	StripeCheckoutSessionID sql.NullString `db:"stripe_checkout_session_id"`
	CheckoutURL             sql.NullString `db:"checkout_url"`
	CurrentPeriodEnd        sql.NullTime   `db:"current_period_end"`
	SucceededAt             sql.NullTime   `db:"succeeded_at"`
	FailedAt                sql.NullTime   `db:"failed_at"`
	FailureReason           sql.NullString `db:"failure_reason"`
	CreatedAt               time.Time      `db:"created_at"`
	UpdatedAt               time.Time      `db:"updated_at"`
}
