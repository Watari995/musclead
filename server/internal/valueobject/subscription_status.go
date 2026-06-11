package valueobject

import "errors"

type SubscriptionStatusCode string

const (
	SubscriptionStatusActive   SubscriptionStatusCode = "active"
	SubscriptionStatusCanceled SubscriptionStatusCode = "canceled"
	SubscriptionStatusExpired  SubscriptionStatusCode = "expired"
)

var ErrInvalidSubscriptionStatus = errors.New("invalid subscription status")

type SubscriptionStatus struct {
	LiteralBase[string]
}

func NewSubscriptionStatusFromString(s string) (*SubscriptionStatus, error) {
	switch SubscriptionStatusCode(s) {
	case SubscriptionStatusActive, SubscriptionStatusCanceled, SubscriptionStatusExpired:
		return &SubscriptionStatus{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidSubscriptionStatus
	}
}

func NewSubscriptionStatusFromCode(c SubscriptionStatusCode) SubscriptionStatus {
	return SubscriptionStatus{LiteralBase: LiteralBase[string]{v: string(c)}}
}
