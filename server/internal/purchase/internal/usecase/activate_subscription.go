package purchaseusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
)

type ActivateSubscription struct {
	subscriptionRepo purchasedomain.SubscriptionRepository
	orderRepo        purchasedomain.SubscriptionOrderRepository
}

func (uc *ActivateSubscription) ActivateSubscription(ctx context.Context, req purchasepublicfunctions.ActivateSubscriptionRequest) error {
	subscription, err := uc.subscriptionRepo.FindByPaymentID(ctx, req.PaymentID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if subscription != nil {
		return nil
	}
	order, err := uc.orderRepo.FindByPaymentID(ctx, req.PaymentID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if order == nil {
		return myerror.NewSubscriptionOrderNotFoundError()
	}
	order.MarkSucceeded(time.Now())
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	orderID := order.ID()
	subscription = purchasedomain.CreateSubscription(req.UserID, req.Plan, &orderID, req.PaymentID, req.ExpiresAt)
	if err := uc.subscriptionRepo.Save(ctx, subscription); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewActivateSubscription(subscriptionRepo purchasedomain.SubscriptionRepository, orderRepo purchasedomain.SubscriptionOrderRepository) *ActivateSubscription {
	return &ActivateSubscription{subscriptionRepo: subscriptionRepo, orderRepo: orderRepo}
}
