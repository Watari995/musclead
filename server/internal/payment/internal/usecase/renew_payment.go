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

// RenewPaymentInput は月次自動更新 Webhook ('invoice.payment_succeeded') の入力。
type RenewPaymentInput struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata
}

// RenewPayment は Stripe の月次自動更新成功時に payments.current_period_end を更新する。
//
// 設計 (ADR 0014, 0018):
//   - Stripe 月次課金成功時に 'invoice.payment_succeeded' Webhook 受信
//   - TX 内で stripe_events Create + payments UPDATE (current_period_end) + payment_events INSERT (renewed) + outbox INSERT (PaymentRenewed)
//   - Pro 期限 (subscriptions.expires_at) は purchase context が outbox を受けて更新
type RenewPayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

// Execute は Webhook 受信時の本処理。 CompletePayment / CancelPayment と同じパターンの TX。
//
// TODO (User 実装):
//   - 詳細は CompletePayment のコメント参照、 status は変えず current_period_end のみ更新
func (uc *RenewPayment) Execute(ctx context.Context, input RenewPaymentInput) error {
	stripeSubscriptionID := input.Payload["subscription"].(string)
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
	payment.MarkRenewed(currentPeriodEnd)

	paymentEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
	}
	paymentEvent := paymentdomain.CreatePaymentEvent(payment.ID(), valueobject.NewPaymentEventTypeFromCode(valueobject.PaymentEventTypeRenewed), paymentEventMetadata)

	outboxEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
		"current_period_end":     currentPeriodEnd.Format(time.RFC3339),
		"subscription_plan":      valueobject.SubscriptionPlanPro,
	}
	outboxEvent := paymentdomain.CreateOutboxEvent(valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypePaymentRenewed), payment.ID().String(), outboxEventMetadata)

	// stripe_events / payments / payment_events / outbox_events を atomic に保存 (ADR 0014, 0018)
	return uc.txManager.Processing(ctx, func(ctx context.Context) error {
		if err := uc.stripeEventRepo.Create(ctx, stripeEvent); err != nil {
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

func NewRenewPayment(
	paymentRepo paymentdomain.PaymentRepository,
	paymentEventRepo paymentdomain.PaymentEventRepository,
	stripeEventRepo paymentdomain.StripeEventRepository,
	outboxEventRepo paymentdomain.OutboxEventRepository,
	txManager dbtx.TransactionManager,
) *RenewPayment {
	return &RenewPayment{
		paymentRepo:      paymentRepo,
		paymentEventRepo: paymentEventRepo,
		stripeEventRepo:  stripeEventRepo,
		outboxEventRepo:  outboxEventRepo,
		txManager:        txManager,
	}
}
