package paymentinfra

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// PaymentEventModel は payment_events テーブルの gorp マッピング。
// migration: sql/migrations/000016_create_payment_events.up.sql
// append-only なので updated_at なし。
type PaymentEventModel struct {
	ID        []byte               `db:"id"`
	PaymentID []byte               `db:"payment_id"`
	EventType string               `db:"event_type"`
	Metadata  valueobject.Metadata `db:"metadata"` // Scanner/Valuer 実装済みで自動シリアライズ
	CreatedAt time.Time            `db:"created_at"`
}
