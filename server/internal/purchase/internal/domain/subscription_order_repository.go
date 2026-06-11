package purchasedomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)
type SubscriptionOrderRepository interface {
	FindPendingByUserID(ctx context.Context, userID valueobject.UserID) (*SubscriptionOrder, error)
	FindByPaymentID(ctx context.Context, paymentID valueobject.PaymentID) (*SubscriptionOrder, error)
	Save(ctx context.Context, order *SubscriptionOrder) error
}
