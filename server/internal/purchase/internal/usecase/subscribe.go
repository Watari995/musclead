package purchaseusecase

import (
	"context"
	"log/slog"

	"github.com/Watari995/musclead/internal/myerror"
	paymentpublicfunctions "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

// SubscribeInput はサブスク申込み時に handler から渡される入力。
// handler は HTTP プロトコルから取得した最小情報 (UserID + Plan) だけを渡す。
type SubscribeInput struct {
	UserID valueobject.UserID
	Plan   valueobject.SubscriptionPlan
}

// SubscribeOutput は client にレスポンスする情報。
type SubscribeOutput struct {
	CheckoutURL valueobject.URL
}

// Subscribe はサブスク申込みのオーケストレータ usecase (ADR 0013)。
//
// 流れ:
//  1. 既存 pending order を user で検索 → あれば再利用、 なければ新規 INSERT
//  2. user の email を userQuery 経由で取得 (cross-module)
//  3. plan → priceID を priceIDByPlan map で解決
//  4. payment.Command.InitiatePayment を呼ぶ (cross-module)
//  5. order に payment_id を紐付け → Save (UPDATE)
//  6. CheckoutURL を Output で返却
//
// 設計メモ:
//   - publicfunctions (paymentCommand / userQuery) は usecase 層で呼ぶ。 handler は触らない (ADR 0013)
//   - business 設定 (priceIDByPlan) は usecase が保持、 handler に漏らさない
//   - 金額は Stripe 側の Price object で管理 (アプリでは保持しない)
//   - subscription (権利状態) はここで作らない。 Webhook 受信時に別途 purchase worker (Phase 9) が作る
type Subscribe struct {
	orderRepo      purchasedomain.SubscriptionOrderRepository
	userQuery      userpublicfunctions.UserQuery
	paymentCommand paymentpublicfunctions.PaymentCommand
	priceIDByPlan  map[valueobject.SubscriptionPlanCode]string
}

func (uc *Subscribe) Execute(ctx context.Context, input SubscribeInput) (SubscribeOutput, error) {
	existing, err := uc.orderRepo.FindPendingByUserID(ctx, input.UserID)
	if err != nil {
		return SubscribeOutput{}, myerror.NewInternalError().Wrap(err)
	}
	// orderがない場合は新しく作成する
	if existing == nil {
		existing = purchasedomain.CreateSubscriptionOrder(input.UserID, input.Plan)
		if err := uc.orderRepo.Save(ctx, existing); err != nil {
			return SubscribeOutput{}, myerror.NewInternalError().Wrap(err)
		}
	}
	output, err := uc.userQuery.GetEmailByUserID(ctx, userpublicfunctions.GetEmailByUserIDInput{UserID: input.UserID})
	if err != nil {
		return SubscribeOutput{}, err // すでにwrap済み
	}
	priceID, ok := uc.priceIDByPlan[input.Plan.Code()]
	if !ok {
		slog.Error("priceID not found", "plan", input.Plan.Code())
		return SubscribeOutput{}, myerror.NewInternalError().SetMessage("priceID not found")
	}

	resp, err := uc.paymentCommand.InitiatePayment(ctx, paymentpublicfunctions.InitiatePaymentRequest{
		UserID:  input.UserID,
		Email:   output.Email,
		PriceID: priceID,
	})
	if err != nil {
		return SubscribeOutput{}, err
	}

	// 既存のorderにpaymentIDを紐づける
	existing.AttachPayment(resp.PaymentID)
	if err := uc.orderRepo.Save(ctx, existing); err != nil {
		return SubscribeOutput{}, myerror.NewInternalError().Wrap(err)
	}

	return SubscribeOutput{CheckoutURL: resp.CheckoutURL}, nil
}

func NewSubscribe(
	orderRepo purchasedomain.SubscriptionOrderRepository,
	paymentCommand paymentpublicfunctions.PaymentCommand,
	userQuery userpublicfunctions.UserQuery,
	priceIDByPlan map[valueobject.SubscriptionPlanCode]string,
) *Subscribe {
	return &Subscribe{
		orderRepo:      orderRepo,
		paymentCommand: paymentCommand,
		userQuery:      userQuery,
		priceIDByPlan:  priceIDByPlan,
	}
}
