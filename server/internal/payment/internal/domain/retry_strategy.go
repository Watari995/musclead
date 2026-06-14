package paymentdomain

import "context"

type RetryStrategy interface {
	OnFailure(ctx context.Context, event *StripeEvent, err error) error
}
