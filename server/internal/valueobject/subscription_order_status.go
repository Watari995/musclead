package valueobject

import "errors"

type SubscriptionOrderStatusCode string

const (
	SubscriptionOrderStatusPending   SubscriptionOrderStatusCode = "pending"
	SubscriptionOrderStatusSucceeded SubscriptionOrderStatusCode = "succeeded"
	SubscriptionOrderStatusFailed    SubscriptionOrderStatusCode = "failed"
)

var ErrInvalidSubscriptionOrderStatus = errors.New("invalid subscription order status")

type SubscriptionOrderStatus struct {
	LiteralBase[string]
}

func NewSubscriptionOrderStatusFromString(s string) (*SubscriptionOrderStatus, error) {
	switch SubscriptionOrderStatusCode(s) {
	case SubscriptionOrderStatusPending, SubscriptionOrderStatusSucceeded, SubscriptionOrderStatusFailed:
		return &SubscriptionOrderStatus{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidSubscriptionOrderStatus
	}
}

func NewSubscriptionOrderStatusFromCode(c SubscriptionOrderStatusCode) SubscriptionOrderStatus {
	return SubscriptionOrderStatus{LiteralBase: LiteralBase[string]{v: string(c)}}
}
