package purchaseusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// GetSubscription は user の現在のサブスク状態 (is_pro / plan / expires_at) を返す Query usecase。
// GET /purchase/subscription のために handler から呼ばれる。 Stripe は叩かず DB のみ。
//
// 流れ:
//  1. subscriptionRepo.FindLatestByUserID で最新 subscription を取得
//  2. 無ければ free (is_pro=false)。 有れば expires_at > NOW() で is_pro を判定 (ADR 0017)
//  3. GetSubscriptionOutput を返す
type GetSubscription struct {
	subscriptionRepo purchasedomain.SubscriptionRepository
}

type GetSubscriptionInput struct {
	UserID valueobject.UserID
}

// GetSubscriptionOutput は handler が JSON 化する。 free のとき Plan / ExpiresAt は nil。
type GetSubscriptionOutput struct {
	IsPro     bool
	Plan      *valueobject.SubscriptionPlan
	ExpiresAt *time.Time
}

func (uc *GetSubscription) Execute(ctx context.Context, input GetSubscriptionInput) (GetSubscriptionOutput, error) {
	subscription, err := uc.subscriptionRepo.FindLatestByUserID(ctx, input.UserID)
	if err != nil {
		return GetSubscriptionOutput{}, myerror.NewInternalError().Wrap(err)
	}
	if subscription == nil {
		return GetSubscriptionOutput{IsPro: false}, nil
	}
	isPro := subscription.IsPro()
	plan := subscription.Plan()
	expiresAt := subscription.ExpiresAt()
	return GetSubscriptionOutput{IsPro: isPro, Plan: &plan, ExpiresAt: &expiresAt}, nil
}

func NewGetSubscription(subscriptionRepo purchasedomain.SubscriptionRepository) *GetSubscription {
	return &GetSubscription{subscriptionRepo: subscriptionRepo}
}
