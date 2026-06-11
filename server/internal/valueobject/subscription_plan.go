package valueobject

import "errors"

type SubscriptionPlanCode string

const (
	SubscriptionPlanPro SubscriptionPlanCode = "pro"
)

var ErrInvalidSubscriptionPlan = errors.New("invalid subscription plan")

type SubscriptionPlan struct {
	LiteralBase[string]
}

func NewSubscriptionPlanFromString(s string) (*SubscriptionPlan, error) {
	switch SubscriptionPlanCode(s) {
	case SubscriptionPlanPro:
		return &SubscriptionPlan{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidSubscriptionPlan
	}
}

func NewSubscriptionPlanFromCode(c SubscriptionPlanCode) SubscriptionPlan {
	return SubscriptionPlan{LiteralBase: LiteralBase[string]{v: string(c)}}
}

func (s SubscriptionPlan) Code() SubscriptionPlanCode {
	return SubscriptionPlanCode(s.Value())
}
