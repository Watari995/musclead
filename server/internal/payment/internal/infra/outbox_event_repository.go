package paymentinfra

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type outboxEventRepository struct {
	dbmap *gorp.DbMap
}

func NewOutboxEventRepository(dbmap *gorp.DbMap) paymentdomain.OutboxEventRepository {
	return &outboxEventRepository{dbmap: dbmap}
}

func (r *outboxEventRepository) FindPending(ctx context.Context, limit int) ([]*paymentdomain.OutboxEvent, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var rows []OutboxEventModel
	_, err := q.Select(&rows,
		`SELECT id, event_type, aggregate_id, payload, published_at, publish_error, created_at, updated_at FROM outbox_events WHERE published_at IS NULL ORDER BY created_at ASC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	outboxEvents := make([]*paymentdomain.OutboxEvent, 0, len(rows))
	for _, row := range rows {
		outboxEvent, err := toOutboxEvent(row)
		if err != nil {
			return nil, err
		}
		outboxEvents = append(outboxEvents, outboxEvent)
	}
	return outboxEvents, nil
}

const upsertOutboxEventSQL = `
INSERT INTO outbox_events (id, event_type, aggregate_id, payload, published_at, publish_error, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    event_type = VALUES(event_type),
    aggregate_id = VALUES(aggregate_id),
    payload = VALUES(payload),
    published_at = VALUES(published_at),
    publish_error = VALUES(publish_error),
    updated_at = VALUES(updated_at)
`

func (r *outboxEventRepository) Save(ctx context.Context, event *paymentdomain.OutboxEvent) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertOutboxEventParams(event)
	if err != nil {
		return err
	}
	if _, err := q.Exec(upsertOutboxEventSQL, params...); err != nil {
		return err
	}
	return nil
}

func toOutboxEvent(row OutboxEventModel) (*paymentdomain.OutboxEvent, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.OutboxEventID](row.ID)
	if err != nil {
		return nil, err
	}
	eventType, err := valueobject.NewOutboxEventTypeFromString(row.EventType)
	if err != nil {
		return nil, err
	}
	aggregateID, err := sqlconv.UUIDStringFromBytes(row.AggregateID)
	if err != nil {
		return nil, err
	}
	publishedAt := sqlconv.FromNullTime(row.PublishedAt)
	publishError := sqlconv.NewStringFromNullString(row.PublishError)
	return paymentdomain.NewOutboxEvent(
		*id,
		*eventType,
		aggregateID,
		row.Payload,
		publishedAt,
		publishError,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func buildUpsertOutboxEventParams(event *paymentdomain.OutboxEvent) ([]any, error) {
	idBytes, err := event.ID().Bytes()
	if err != nil {
		return nil, err
	}
	eventType := event.EventType().Value()
	aggregateIDBytes, err := sqlconv.UUIDStringToBytes(event.AggregateID())
	if err != nil {
		return nil, err
	}
	payload := event.Payload()
	publishedAt := sqlconv.ToNullTime(event.PublishedAt())
	publishError := sqlconv.StringPtrToNullString(event.PublishError())
	createdAt := event.CreatedAt()
	updatedAt := event.UpdatedAt()
	return []any{idBytes, eventType, aggregateIDBytes, payload, publishedAt, publishError, createdAt, updatedAt}, nil
}
