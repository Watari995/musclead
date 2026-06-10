package paymentdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// StripeEvent は Stripe Webhook 受信記録 + 冪等性キーを保持する。
//
// 設計 (ADR 0014):
//   - 受信した全 Stripe event を生 payload で残す (監査 / デバッグ / 再処理用)
//   - stripe_event_id (Stripe 側の evt_xxx) を UNIQUE 制約にして二重処理を物理的に防ぐ
//   - 処理完了で processed_at を SET、 失敗時は processing_error を SET
//
// migration: sql/migrations/000017_create_stripe_events.up.sql
type StripeEvent struct {
	id              valueobject.StripeEventID
	stripeEventID   string
	eventType       string
	payload         valueobject.Metadata
	processedAt     *time.Time
	processingError *string
	createdAt       time.Time
	updatedAt       time.Time
}

func (s *StripeEvent) ID() valueobject.StripeEventID {
	return s.id
}

func (s *StripeEvent) StripeEventID() string {
	return s.stripeEventID
}

func (s *StripeEvent) EventType() string {
	return s.eventType
}

func (s *StripeEvent) Payload() valueobject.Metadata {
	return s.payload
}

func (s *StripeEvent) ProcessedAt() *time.Time {
	return s.processedAt
}

func (s *StripeEvent) IsProcessed() bool {
	return s.processedAt != nil
}

func (s *StripeEvent) MarkProcessed() {
	if s.IsProcessed() {
		return
	}
	now := time.Now()
	s.processedAt = &now
	s.updatedAt = now
}

func (s *StripeEvent) ProcessingError() *string {
	return s.processingError
}

// 最新のエラーで上書きする
func (s *StripeEvent) MarkFailed(err string) {
	s.processingError = &err
	s.updatedAt = time.Now()
}

func (s *StripeEvent) CreatedAt() time.Time {
	return s.createdAt
}

func (s *StripeEvent) UpdatedAt() time.Time {
	return s.updatedAt
}

func CreateStripeEvent(
	stripeEventID string,
	eventType string,
	payload valueobject.Metadata,
) *StripeEvent {
	return &StripeEvent{
		id:              valueobject.NewPrimaryID[valueobject.StripeEventID](),
		stripeEventID:   stripeEventID,
		eventType:       eventType,
		payload:         payload,
		processedAt:     nil,
		processingError: nil,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
	}
}

func NewStripeEvent(
	id valueobject.StripeEventID,
	stripeEventID string,
	eventType string,
	payload valueobject.Metadata,
	processedAt *time.Time,
	processingError *string,
	createdAt time.Time,
	updatedAt time.Time,
) *StripeEvent {
	return &StripeEvent{
		id:              id,
		stripeEventID:   stripeEventID,
		eventType:       eventType,
		payload:         payload,
		processedAt:     processedAt,
		processingError: processingError,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}
