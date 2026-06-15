package purchaseusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
)

// CancelSubscription は解約 (Stripe 'customer.subscription.deleted') を受けて subscription を終了状態にする。
// billing handler が CancelPayment の後に呼ぶ。 引き当ては CancelPayment が返す PaymentID で行う。
//
// 流れ:
//  1. subscriptionRepo.FindByPaymentID で対象 subscription を取得
//  2. sub.MarkExpired() で終了状態にする (MarkCanceled との使い分けは ADR 0017 を参照して決める)
//  3. subscriptionRepo.Save で永続化
//
// 手本: renew_subscription.go / activate_subscription.go
type CancelSubscription struct {
	subscriptionRepo purchasedomain.SubscriptionRepository
}

func (uc *CancelSubscription) CancelSubscription(ctx context.Context, req purchasepublicfunctions.CancelSubscriptionRequest) error {
	subscription, err := uc.subscriptionRepo.FindByPaymentID(ctx, req.PaymentID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if subscription == nil {
		return myerror.NewSubscriptionNotFoundError()
	}
	subscription.MarkExpired()
	if err := uc.subscriptionRepo.Save(ctx, subscription); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewCancelSubscription(subscriptionRepo purchasedomain.SubscriptionRepository) *CancelSubscription {
	return &CancelSubscription{subscriptionRepo: subscriptionRepo}
}
