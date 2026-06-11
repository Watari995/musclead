package purchaseusecase

import (
	"context"

	paymentpublic "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasedomain "github.com/Watari995/musclead/internal/purchase/internal/domain"
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

// UserQuery は user 集約の参照用 interface。
// user モジュールが publicfunctions.UserQuery として公開する想定 (Phase 後半で追加)。
// 注入することで purchase usecase が user の email を取得できる。
type UserQuery interface {
	GetEmailByUserID(ctx context.Context, userID valueobject.UserID) (valueobject.Email, error)
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
//   - business 設定 (priceIDByPlan / proAmount) も usecase が保持、 handler に漏らさない
//   - subscription (権利状態) はここで作らない。 Webhook 受信時に別途 purchase worker (Phase 9) が作る
type Subscribe struct {
	orderRepo      purchasedomain.SubscriptionOrderRepository
	paymentCommand paymentpublic.PaymentCommand
	userQuery      UserQuery
	priceIDByPlan  map[valueobject.SubscriptionPlanCode]string
	proAmount      valueobject.NonNegativeInt
}

func (uc *Subscribe) Execute(ctx context.Context, input SubscribeInput) (SubscribeOutput, error) {
	// TODO (User 実装): 上記の「流れ」 を実装
	return SubscribeOutput{}, nil
}

func NewSubscribe(
	orderRepo purchasedomain.SubscriptionOrderRepository,
	paymentCommand paymentpublic.PaymentCommand,
	userQuery UserQuery,
	priceIDByPlan map[valueobject.SubscriptionPlanCode]string,
	proAmount valueobject.NonNegativeInt,
) *Subscribe {
	return &Subscribe{
		orderRepo:      orderRepo,
		paymentCommand: paymentCommand,
		userQuery:      userQuery,
		priceIDByPlan:  priceIDByPlan,
		proAmount:      proAmount,
	}
}
