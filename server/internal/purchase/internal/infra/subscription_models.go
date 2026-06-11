package purchaseinfra

import (
	"database/sql"
	"time"
)

// SubscriptionModel は subscriptions テーブルの gorp マッピング。
// migration: sql/migrations/000021_create_subscriptions.up.sql
type SubscriptionModel struct {
	ID                  []byte       `db:"id"`
	UserID              []byte       `db:"user_id"`
	Plan                string       `db:"plan"`
	Status              string       `db:"status"`
	SubscriptionOrderID []byte       `db:"subscription_order_id"` // nullable: admin 手動作成では NULL
	PaymentID           []byte       `db:"payment_id"`            // NOT NULL: Webhook で INSERT する時点で必ず存在
	ActivatedAt         time.Time    `db:"activated_at"`
	ExpiresAt           time.Time    `db:"expires_at"`
	CanceledAt          sql.NullTime `db:"canceled_at"`
	CreatedAt           time.Time    `db:"created_at"`
	UpdatedAt           time.Time    `db:"updated_at"`
}
