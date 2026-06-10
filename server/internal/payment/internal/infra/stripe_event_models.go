package paymentinfra

import (
	"database/sql"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// StripeEventModel は stripe_events テーブルの gorp マッピング。
// migration: sql/migrations/000017_create_stripe_events.up.sql
// stripe_event_id (Stripe 側の evt_xxx) は UNIQUE 制約付き。
type StripeEventModel struct {
	ID              []byte               `db:"id"`
	StripeEventID   string               `db:"stripe_event_id"`
	EventType       string               `db:"event_type"`
	Payload         valueobject.Metadata `db:"payload"` // JSON、 Scanner/Valuer 実装済み
	ProcessedAt     sql.NullTime         `db:"processed_at"`
	ProcessingError sql.NullString       `db:"processing_error"`
	CreatedAt       time.Time            `db:"created_at"`
	UpdatedAt       time.Time            `db:"updated_at"`
}
