package outboxinfra

import (
	"context"
	"fmt"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlquery"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type outboxEventRepository struct {
	dbmap *gorp.DbMap
}

func NewOutboxEventRepository(dbmap *gorp.DbMap) shareddomain.OutboxEventRepository {
	return &outboxEventRepository{dbmap: dbmap}
}

func buildFindPendingByEventTypesSQL(eventTypes []valueobject.OutboxEventType, limit int) (string, []any) {
	values := make([]string, 0, len(eventTypes))
	for _, v := range eventTypes {
		value := v.Value()
		values = append(values, value)
	}
	placeholders, args := sqlquery.InPlaceholders(values)

	return fmt.Sprintf(`
	SELECT id, event_type, aggregate_id, payload, published_at, publish_error, created_at, updated_at
	FROM outbox_events
	WHERE published_at IS NULL
	AND event_type IN (%s)
	ORDER BY created_at ASC LIMIT %d
	`, placeholders, limit), args
}

func (r *outboxEventRepository) FindPendingByEventTypes(ctx context.Context, eventTypes []valueobject.OutboxEventType, limit int) ([]*shareddomain.OutboxEvent, error) {
	if len(eventTypes) == 0 {
		return nil, nil // 即時終了
	}
	q := dbtx.Querier(ctx, r.dbmap)

	sqlStr, args := buildFindPendingByEventTypesSQL(eventTypes, limit)

	var rows []OutboxEventModel
	_, err := q.Select(&rows, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	outboxEvents := make([]*shareddomain.OutboxEvent, 0, len(rows))
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

func (r *outboxEventRepository) Save(ctx context.Context, event *shareddomain.OutboxEvent) error {
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

func toOutboxEvent(row OutboxEventModel) (*shareddomain.OutboxEvent, error) {
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
	return shareddomain.NewOutboxEvent(
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

func buildUpsertOutboxEventParams(event *shareddomain.OutboxEvent) ([]any, error) {
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
