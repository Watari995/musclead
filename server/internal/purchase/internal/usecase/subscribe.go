package purchaseusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// SubscribeInput はサブスク申込み時に handler から渡される入力。
type SubscribeInput struct {
	UserID  valueobject.UserID
	Email   valueobject.Email
	Amount  valueobject.NonNegativeInt
	PriceID string                       // Stripe Price ID (plan に応じて handler / main.go で解決)
	Plan    valueobject.SubscriptionPlan // pro / 将来 pro_annual 等
}

// SubscribeOutput は client にレスポンスする情報。
type SubscribeOutput struct {
	CheckoutURL valueobject.URL
}

// Subscribe はサブスク申込みのオーケストレータ usecase (ADR 0013)。
//
// 流れ:
//  1. 既存 pending order を user で検索 → あれば再利用、 なければ CreateSubscriptionOrder (pending、 plan は入力で受け取る)
//  2. paymentCommand.InitiatePayment(...) を呼ぶ (Phase 1 で公開済み)
//     → payment 集約が paymentID と Stripe Checkout Session URL を返す
//  3. order に payment_id を紐付ける (AttachPayment) → Save (UPDATE)
//  4. CheckoutURL を client に返却
//
// 設計メモ:
//   - plan は input で受け取って汎用化 (将来 pro_annual 等を追加する時 interface を変えなくて済む)
//   - subscription (権利状態) はここでは作らない。 Webhook 受信時に別途 purchase worker (Phase 9) が作る
//   - payment 側の Stripe API 呼び出しは payment 集約に閉じる (purchase は Stripe を知らない)
type Subscribe struct {
	orderRepo      purchasedomain.SubscriptionOrderRepository
	paymentCommand publicfunctions.PaymentCommand
}

func (uc *Subscribe) Execute(ctx context.Context, input SubscribeInput) (SubscribeOutput, error) {
	order, err := uc.orderRepo.FindPendingByUserID(ctx, input.UserID)
	if err != nil {
		return SubscribeOutput{}, err
	}
	if order == nil {
		order = purchasedomain.CreateSubscriptionOrder(input.UserID, input.Plan)
		if err := uc.orderRepo.Save(ctx, order); err != nil {
			return SubscribeOutput{}, err
		}
	}
	resp, err := uc.paymentCommand.InitiatePayment(ctx, publicfunctions.InitiatePaymentRequest{
		UserID:  input.UserID,
		Email:   input.Email,
		Amount:  input.Amount,
		PriceID: input.PriceID,
	})
	if err != nil {
		return SubscribeOutput{}, err
	}
	order.AttachPayment(resp.PaymentID)
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return SubscribeOutput{}, err
	}
	return SubscribeOutput{CheckoutURL: resp.CheckoutURL}, nil
}

func NewSubscribe(orderRepo purchasedomain.SubscriptionOrderRepository, paymentCommand publicfunctions.PaymentCommand) *Subscribe {
	return &Subscribe{
		orderRepo:      orderRepo,
		paymentCommand: paymentCommand,
	}
}
