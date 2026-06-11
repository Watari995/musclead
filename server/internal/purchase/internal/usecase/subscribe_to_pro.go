package purchaseusecase

import (
	"context"
	"errors"

	paymentpublic "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// SubscribeToProInput は purchase handler から渡される入力。
type SubscribeToProInput struct {
	UserID  valueobject.UserID
	Email   valueobject.Email
	Amount  valueobject.NonNegativeInt // 480 (税込 JPY)
	PriceID string                     // Stripe Price ID (商品差分、 env 経由)
}

// SubscribeToProOutput は client にレスポンスする情報。
type SubscribeToProOutput struct {
	CheckoutURL valueobject.URL
}

// SubscribeToPro は Pro 申込みのオーケストレータ usecase (ADR 0013)。
//
// 流れ:
//  1. 既存 pending order を user で検索 → あれば再利用、 なければ CreateSubscriptionOrder (pending)
//  2. paymentCommand.InitiatePayment(...) を呼ぶ (Phase 1 で公開済み)
//     → payment 集約が paymentID と Stripe Checkout Session URL を返す
//  3. order に payment_id を紐付ける (AttachPayment) → Save (UPDATE)
//  4. CheckoutURL を client に返却
//
// 設計メモ:
//   - subscription (権利状態) はここでは作らない。 Webhook 受信時に別途 purchase worker (Phase 9) が作る
//   - payment 側の Stripe API 呼び出しは payment 集約に閉じる (purchase は Stripe を知らない)
type SubscribeToPro struct {
	orderRepo      purchasedomain.SubscriptionOrderRepository
	paymentCommand paymentpublic.PaymentCommand
}

// TODO (User 実装):
//
//	// 1. 既存 pending order を user で検索
//	order, err := uc.orderRepo.FindPendingByUserID(ctx, input.UserID)
//	if err != nil { return SubscribeToProOutput{}, err }
//	if order == nil {
//	    order = purchasedomain.CreateSubscriptionOrder(input.UserID, valueobject.NewSubscriptionPlanFromCode(valueobject.SubscriptionPlanPro))
//	    if err := uc.orderRepo.Save(ctx, order); err != nil { return SubscribeToProOutput{}, err }
//	}
//
//	// 2. payment.InitiatePayment を呼ぶ
//	resp, err := uc.paymentCommand.InitiatePayment(ctx, paymentpublic.InitiatePaymentRequest{
//	    UserID:  input.UserID,
//	    Email:   input.Email,
//	    Amount:  input.Amount,
//	    PriceID: input.PriceID,
//	})
//	if err != nil { return SubscribeToProOutput{}, err }
//
//	// 3. order に payment_id を紐付け
//	order.AttachPayment(resp.PaymentID)
//	if err := uc.orderRepo.Save(ctx, order); err != nil { return SubscribeToProOutput{}, err }
//
//	// 4. CheckoutURL 返却
//	return SubscribeToProOutput{CheckoutURL: resp.CheckoutURL}, nil
func (uc *SubscribeToPro) Execute(ctx context.Context, input SubscribeToProInput) (SubscribeToProOutput, error) {
	return SubscribeToProOutput{}, errors.New("SubscribeToPro: not implemented")
}

func NewSubscribeToPro(orderRepo purchasedomain.SubscriptionOrderRepository, paymentCommand paymentpublic.PaymentCommand) *SubscribeToPro {
	return &SubscribeToPro{
		orderRepo:      orderRepo,
		paymentCommand: paymentCommand,
	}
}
