package valueobject

import "errors"

type PaymentEventTypeCode string

const (
	PaymentEventTypeInitiated PaymentEventTypeCode = "initiated"
	PaymentEventTypeSucceeded PaymentEventTypeCode = "succeeded"
	PaymentEventTypeFailed    PaymentEventTypeCode = "failed"
	PaymentEventTypeCanceled  PaymentEventTypeCode = "canceled"
	PaymentEventTypeRenewed   PaymentEventTypeCode = "renewed"
)

var ErrInvalidPaymentEventType = errors.New("invalid payment event type")

type PaymentEventType struct {
	LiteralBase[string]
}

func NewPaymentEventTypeFromString(s string) (*PaymentEventType, error) {
	switch PaymentEventTypeCode(s) {
	case PaymentEventTypeInitiated, PaymentEventTypeSucceeded, PaymentEventTypeFailed, PaymentEventTypeCanceled, PaymentEventTypeRenewed:
		return &PaymentEventType{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidPaymentEventType
	}
}
