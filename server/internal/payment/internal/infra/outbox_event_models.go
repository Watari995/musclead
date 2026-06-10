package paymentinfra

import (
	"database/sql"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// OutboxEventModel は outbox_events テーブルの gorp マッピング。
// migration: sql/migrations/000018_create_outbox_events.up.sql
// Webhook 処理の TX 内で INSERT、 TX 外で SNS publish。
type OutboxEventModel struct {
	ID           []byte               `db:"id"`
	EventType    string               `db:"event_type"`
	AggregateID  []byte               `db:"aggregate_id"` // payment_id (BINARY(16))
	Payload      valueobject.Metadata `db:"payload"`      // JSON、 Scanner/Valuer 実装済み
	PublishedAt  sql.NullTime         `db:"published_at"`
	PublishError sql.NullString       `db:"publish_error"`
	CreatedAt    time.Time            `db:"created_at"`
	UpdatedAt    time.Time            `db:"updated_at"`
}
