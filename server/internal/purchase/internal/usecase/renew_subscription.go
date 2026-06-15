package purchaseusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
)

type RenewSubscription struct {
	subscriptionRepo purchasedomain.SubscriptionRepository
}

func (uc *RenewSubscription) RenewSubscription(ctx context.Context, req purchasepublicfunctions.RenewSubscriptionRequest) error {
	subscription, err := uc.subscriptionRepo.FindByPaymentID(ctx, req.PaymentID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if subscription == nil {
		return myerror.NewSubscriptionNotFoundError()
	}
	subscription.Renew(req.ExpiresAt)
	if err := uc.subscriptionRepo.Save(ctx, subscription); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewRenewSubscription(subscriptionRepo purchasedomain.SubscriptionRepository) *RenewSubscription {
	return &RenewSubscription{subscriptionRepo: subscriptionRepo}
}
