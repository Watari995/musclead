package paymentusecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

// CompletePaymentInput は Webhook 'checkout.session.completed' 受信時に handler から渡される。
type CompletePaymentInput struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata // event.Data の生 JSON を Metadata に詰めたもの
}

// CompletePayment は Stripe Checkout 完了時に payment を succeeded に遷移させる。
//
// 設計 (ADR 0014, 0018):
//   - TX 内で全部 atomic に実行
//   - stripe_events Create (UNIQUE 違反は ErrStripeEventAlreadyExists で no-op = 冪等性吸収)
//   - payments UPDATE (succeeded、 stripe_subscription_id, current_period_end, succeeded_at)
//   - payment_events INSERT (succeeded)
//   - outbox_events INSERT (PaymentSucceeded、 email worker 用)
type CompletePayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

func (uc *CompletePayment) Execute(ctx context.Context, input CompletePaymentInput) error {
	stripeSubscriptionID, ok := input.Payload["subscription"].(string)
	if !ok {
		return fmt.Errorf("subscription is not a string")
	}
	periodEndRaw, ok := input.Payload["current_period_end"].(float64)
	if !ok {
		return fmt.Errorf("current_period_end is not a float64")
	}
	currentPeriodEnd := time.Unix(int64(periodEndRaw), 0)

	stripeEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
		"current_period_end":     currentPeriodEnd.Format(time.RFC3339),
	}
	stripeEvent := paymentdomain.CreateStripeEvent(input.StripeEventID, input.EventType, stripeEventMetadata)

	payment, err := uc.paymentRepo.FindByStripeSubscriptionID(ctx, stripeSubscriptionID)
	if err != nil {
		return err
	}
	if payment == nil {
		return paymentdomain.ErrPaymentNotFound
	}
	payment.MarkSucceeded(stripeSubscriptionID, currentPeriodEnd)

	paymentEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
		"subscription_plan":      valueobject.SubscriptionPlanPro,
	}
	paymentEvent := paymentdomain.CreatePaymentEvent(payment.ID(), valueobject.NewPaymentEventTypeFromCode(valueobject.PaymentEventTypeSucceeded), paymentEventMetadata)

	outboxEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
		"current_period_end":     currentPeriodEnd.Format(time.RFC3339),
		"subscription_plan":      valueobject.SubscriptionPlanPro,
	}
	outboxEvent := paymentdomain.CreateOutboxEvent(valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypePaymentSucceeded), payment.ID().String(), outboxEventMetadata)

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

func NewCompletePayment(
	paymentRepo paymentdomain.PaymentRepository,
	paymentEventRepo paymentdomain.PaymentEventRepository,
	stripeEventRepo paymentdomain.StripeEventRepository,
	outboxEventRepo paymentdomain.OutboxEventRepository,
	txManager dbtx.TransactionManager,
) *CompletePayment {
	return &CompletePayment{
		paymentRepo:      paymentRepo,
		paymentEventRepo: paymentEventRepo,
		stripeEventRepo:  stripeEventRepo,
		outboxEventRepo:  outboxEventRepo,
		txManager:        txManager,
	}
}
