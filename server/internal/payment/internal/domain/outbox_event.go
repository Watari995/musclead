package paymentdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// OutboxEvent は payment context が発行する通知用 outbox (ADR 0015)。
//
// 設計:
//   - Webhook 処理の TX 内で INSERT し、 TX 外で SNS publish する
//   - 即時 publish 成功時は published_at を SET
//   - 即時 publish 失敗時は 1 分後の outbox-relay Lambda が拾う (failsafe)
//
// migration: sql/migrations/000018_create_outbox_events.up.sql
//
// aggregateID は string (将来 payment 以外の集約 ID も入る前提で汎用化)。
type OutboxEvent struct {
	id           valueobject.OutboxEventID
	eventType    valueobject.OutboxEventType
	aggregateID  string
	payload      valueobject.Metadata
	publishedAt  *time.Time
	publishError *string
	createdAt    time.Time
	updatedAt    time.Time
}

func (o *OutboxEvent) ID() valueobject.OutboxEventID {
	return o.id
}

func (o *OutboxEvent) EventType() valueobject.OutboxEventType {
	return o.eventType
}

func (o *OutboxEvent) AggregateID() string {
	return o.aggregateID
}

func (o *OutboxEvent) Payload() valueobject.Metadata {
	return o.payload
}

func (o *OutboxEvent) PublishedAt() *time.Time {
	return o.publishedAt
}

func (o *OutboxEvent) IsPublished() bool {
	return o.publishedAt != nil
}

func (o *OutboxEvent) MarkPublished() {
	if o.IsPublished() {
		return
	}
	now := time.Now()
	o.publishedAt = &now
	o.updatedAt = now
}

func (o *OutboxEvent) PublishError() *string {
	return o.publishError
}

func (o *OutboxEvent) MarkPublishFailed(err string) {
	o.publishError = &err
	o.updatedAt = time.Now()
}

func (o *OutboxEvent) CreatedAt() time.Time {
	return o.createdAt
}

func (o *OutboxEvent) UpdatedAt() time.Time {
	return o.updatedAt
}

func CreateOutboxEvent(
	eventType valueobject.OutboxEventType,
	aggregateID string,
	payload valueobject.Metadata,
) *OutboxEvent {
	return &OutboxEvent{
		id:           valueobject.NewPrimaryID[valueobject.OutboxEventID](),
		eventType:    eventType,
		aggregateID:  aggregateID,
		payload:      payload,
		publishedAt:  nil,
		publishError: nil,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}
}

func NewOutboxEvent(
	id valueobject.OutboxEventID,
	eventType valueobject.OutboxEventType,
	aggregateID string,
	payload valueobject.Metadata,
	publishedAt *time.Time,
	publishError *string,
	createdAt time.Time,
	updatedAt time.Time,
) *OutboxEvent {
	return &OutboxEvent{
		id:           id,
		eventType:    eventType,
		aggregateID:  aggregateID,
		payload:      payload,
		publishedAt:  publishedAt,
		publishError: publishError,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
