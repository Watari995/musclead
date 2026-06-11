package paymentusecase

import (
	"context"
	"errors"
	"fmt"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

// CancelPaymentInput は Webhook 'customer.subscription.deleted' 受信時の入力。
type CancelPaymentInput struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata
}

// CancelPayment は Stripe 解約完了時に payment を canceled に遷移させる。
//
// 設計 (ADR 0014, 0017, 0018):
//   - Customer Portal で解約 → Stripe が cancel_at_period_end = true → 期末で削除 → 本 usecase
//   - TX 内で全部 atomic に実行
//   - stripe_events Create (冪等性)
//   - payments UPDATE (canceled)
//   - payment_events INSERT (canceled)
//   - outbox_events INSERT (PaymentCanceled、 email 用)
type CancelPayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

// Execute は Webhook 受信時の本処理。 CompletePayment と同じパターンの TX。
//
// TODO (User 実装):
//   - txManager.Processing 内で stripe_events Create + payments UPDATE + payment_events + outbox INSERT
//   - 詳細は CompletePayment のコメント参照
func (uc *CancelPayment) Execute(ctx context.Context, input CancelPaymentInput) error {
	stripeSubscriptionID, ok := input.Payload["subscription"].(string)
	if !ok {
		return fmt.Errorf("subscription is not a string")
	}
	stripeEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
	}
	stripeEvent := paymentdomain.CreateStripeEvent(input.StripeEventID, input.EventType, stripeEventMetadata)
	payment, err := uc.paymentRepo.FindByStripeSubscriptionID(ctx, stripeSubscriptionID)
	if err != nil {
		return err
	}
	if payment == nil {
		return paymentdomain.ErrPaymentNotFound
	}
	payment.MarkCanceled()

	paymentEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
	}
	paymentEvent := paymentdomain.CreatePaymentEvent(payment.ID(), valueobject.NewPaymentEventTypeFromCode(valueobject.PaymentEventTypeCanceled), paymentEventMetadata)

	outboxEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
		"subscription_plan":      valueobject.SubscriptionPlanPro,
	}
	outboxEvent := paymentdomain.CreateOutboxEvent(valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypePaymentCanceled), payment.ID().String(), outboxEventMetadata)

	// stripe_events / payments / payment_events / outbox_events を atomic に保存 (ADR 0014, 0018)
	return uc.txManager.Processing(ctx, func(ctx context.Context) error {
		if err := uc.stripeEventRepo.Create(ctx, stripeEvent); err != nil {
			// UNIQUE 違反 = 既に処理済みの Webhook 重複受信、 no-op で正常終了 (冪等性吸収)
			if errors.Is(err, paymentdomain.ErrStripeEventAlreadyExists) {
				return nil
			}
			return err
		}
		if err := uc.paymentRepo.Save(ctx, payment); err != nil {
			return err
		}
		if err := uc.paymentEventRepo.Create(ctx, paymentEvent); err != nil {
			return err
		}
		if err := uc.outboxEventRepo.Save(ctx, outboxEvent); err != nil {
			return err
		}
		return nil
	})
}

func NewCancelPayment(
	paymentRepo paymentdomain.PaymentRepository,
	paymentEventRepo paymentdomain.PaymentEventRepository,
	stripeEventRepo paymentdomain.StripeEventRepository,
	outboxEventRepo paymentdomain.OutboxEventRepository,
	txManager dbtx.TransactionManager,
) *CancelPayment {
	return &CancelPayment{
		paymentRepo:      paymentRepo,
		paymentEventRepo: paymentEventRepo,
		stripeEventRepo:  stripeEventRepo,
		outboxEventRepo:  outboxEventRepo,
		txManager:        txManager,
	}
}
