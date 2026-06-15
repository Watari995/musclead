package valueobject

import "errors"

type OutboxEventTypeCode string

const (
	OutboxEventTypePaymentSucceeded OutboxEventTypeCode = "PaymentSucceeded"
	OutboxEventTypePaymentFailed    OutboxEventTypeCode = "PaymentFailed"
	OutboxEventTypePaymentCanceled  OutboxEventTypeCode = "PaymentCanceled"
	OutboxEventTypePaymentRenewed   OutboxEventTypeCode = "PaymentRenewed"
)

var ErrInvalidOutboxEventType = errors.New("invalid outbox event type")

type OutboxEventType struct {
	LiteralBase[string]
}

func NewOutboxEventTypeFromString(s string) (*OutboxEventType, error) {
	switch OutboxEventTypeCode(s) {
	case OutboxEventTypePaymentSucceeded, OutboxEventTypePaymentFailed, OutboxEventTypePaymentCanceled, OutboxEventTypePaymentRenewed:
		return &OutboxEventType{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidOutboxEventType
	}
}

func NewOutboxEventTypeFromCode(c OutboxEventTypeCode) OutboxEventType {
	return OutboxEventType{LiteralBase: LiteralBase[string]{v: string(c)}}
}

func (o OutboxEventType) IsPaymentSucceeded() bool {
	return o.Value() == string(OutboxEventTypePaymentSucceeded)
}
