package purchaseinfra

import (
	"database/sql"
	"time"
)

// SubscriptionOrderModel は subscription_orders テーブルの gorp マッピング。
// migration: sql/migrations/000020_create_subscription_orders.up.sql
type SubscriptionOrderModel struct {
	ID          []byte         `db:"id"`
	UserID      []byte         `db:"user_id"`
	Plan        string         `db:"plan"`
	Status      string         `db:"status"`
	PaymentID   []byte         `db:"payment_id"` // nullable: 申込開始時は NULL、 payment 発行後に SET
	SucceededAt sql.NullTime   `db:"succeeded_at"`
	FailedAt    sql.NullTime   `db:"failed_at"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}
