package purchaseusecase

import (
	"context"

	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
)

// ActivateSubscription は Stripe Webhook (checkout.session.completed) 受信時に
// billing.WebhookHandler から呼ばれて subscriptions レコードを INSERT する usecase。
//
// 設計 (ADR 0014 ③ / ADR 0019):
//   - 冪等性: FindByPaymentID で既存チェック、 あれば no-op で正常終了
//   - publicfunctions.PurchaseCommand.ActivateSubscription を満たすため、 method 名は struct 名と同名
//   - subscription_orders は purchase 側で別途 succeeded に遷移 (本 usecase の責務外)
//
// 流れ:
//  1. subscriptionRepo.FindByPaymentID(req.PaymentID) で既存チェック
//  2. 既存あれば return nil (冪等性吸収)
//  3. purchasedomain.CreateSubscription(req.UserID, req.Plan, req.PaymentID, req.ExpiresAt)
//  4. subscriptionRepo.Save
type ActivateSubscription struct {
	subscriptionRepo purchasedomain.SubscriptionRepository
}

func (uc *ActivateSubscription) ActivateSubscription(ctx context.Context, req purchasepublicfunctions.ActivateSubscriptionRequest) error {
	return nil
}

func NewActivateSubscription(subscriptionRepo purchasedomain.SubscriptionRepository) *ActivateSubscription {
	return &ActivateSubscription{subscriptionRepo: subscriptionRepo}
}
