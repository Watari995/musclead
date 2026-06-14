package paymentinfra

import (
	"context"
	"database/sql"
	"errors"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlerr"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type stripeEventRepository struct {
	dbmap *gorp.DbMap
}

func NewStripeEventRepository(dbmap *gorp.DbMap) paymentdomain.StripeEventRepository {
	return &stripeEventRepository{dbmap: dbmap}
}

func (r *stripeEventRepository) FindByStripeEventID(ctx context.Context, stripeEventID string) (*paymentdomain.StripeEvent, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var row StripeEventModel
	err := q.SelectOne(&row,
		`SELECT id, stripe_event_id, event_type, payload, processed_at, processing_error, created_at, updated_at FROM stripe_events WHERE stripe_event_id = ?`, stripeEventID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toStripeEvent(row)
}

const insertStripeEventSQL = `
INSERT INTO stripe_events (id, stripe_event_id, event_type, payload, processed_at, processing_error, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

func (r *stripeEventRepository) Create(ctx context.Context, event *paymentdomain.StripeEvent) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildInsertStripeEventParams(event)
	if err != nil {
		return err
	}
	if _, err := q.Exec(insertStripeEventSQL, params...); err != nil {
		if sqlerr.IsDuplicateKey(err) {
			return paymentdomain.ErrStripeEventAlreadyExists
		}
		return err
	}
	return nil
}

const updateStripeEventProcessStatusSQL = `
UPDATE stripe_events SET processed_at = ?, processing_error = ?, updated_at = ? WHERE stripe_event_id = ?
`

func (r *stripeEventRepository) UpdateProcessStatus(ctx context.Context, event *paymentdomain.StripeEvent) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpdateStripeEventProcessStatusParams(event)
	if err != nil {
		return err
	}
	if _, err := q.Exec(updateStripeEventProcessStatusSQL, params...); err != nil {
		return err
	}
	return nil
}

func toStripeEvent(row StripeEventModel) (*paymentdomain.StripeEvent, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.StripeEventID](row.ID)
	if err != nil {
		return nil, err
	}
	processedAt := sqlconv.FromNullTime(row.ProcessedAt)
	processingError := sqlconv.NewStringFromNullString(row.ProcessingError)
	return paymentdomain.NewStripeEvent(
		*id,
		row.StripeEventID,
		row.EventType,
		row.Payload,
		processedAt,
		processingError,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func buildInsertStripeEventParams(event *paymentdomain.StripeEvent) ([]any, error) {
	bytes, err := event.ID().Bytes()
	if err != nil {
		return nil, err
	}
	processedAt := sqlconv.ToNullTime(event.ProcessedAt())
	processingError := sqlconv.StringPtrToNullString(event.ProcessingError())
	return []any{bytes, event.StripeEventID(), event.EventType(), event.Payload(), processedAt, processingError, event.CreatedAt(), event.UpdatedAt()}, nil
}

func buildUpdateStripeEventProcessStatusParams(event *paymentdomain.StripeEvent) ([]any, error) {
	processedAt := sqlconv.ToNullTime(event.ProcessedAt())
	processingError := sqlconv.StringPtrToNullString(event.ProcessingError())
	return []any{processedAt, processingError, event.UpdatedAt(), event.StripeEventID()}, nil
}
