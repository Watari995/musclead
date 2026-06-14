package paymentinfra

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/go-gorp/gorp/v3"
)

type paymentEventRepository struct {
	dbmap *gorp.DbMap
}

func NewPaymentEventRepository(dbmap *gorp.DbMap) paymentdomain.PaymentEventRepository {
	return &paymentEventRepository{dbmap: dbmap}
}

const insertPaymentEventSQL = `
INSERT INTO payment_events (id, payment_id, event_type, metadata, created_at)
VALUES (?, ?, ?, ?, ?)
`

func (r *paymentEventRepository) Create(ctx context.Context, event *paymentdomain.PaymentEvent) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildInsertPaymentEventParams(event)
	if err != nil {
		return err
	}
	_, err = q.Exec(insertPaymentEventSQL, params...)
	if err != nil {
		return err
	}
	return nil
}

func buildInsertPaymentEventParams(event *paymentdomain.PaymentEvent) ([]any, error) {
	bytes, err := event.ID().Bytes()
	if err != nil {
		return nil, err
	}
	paymentIDBytes, err := event.PaymentID().Bytes()
	if err != nil {
		return nil, err
	}
	return []any{
		bytes,
		paymentIDBytes,
		event.EventType().Value(),
		event.Metadata(),
		event.CreatedAt(),
	}, nil
}
