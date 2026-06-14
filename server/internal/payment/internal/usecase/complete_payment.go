package paymentusecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CompletePayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

func (uc *CompletePayment) CompletePayment(ctx context.Context, input publicfunctions.CompletePaymentRequest) (publicfunctions.CompletePaymentResponse, error) {
	stripeSubscriptionID, ok := input.Payload["subscription"].(string)
	if !ok {
		return publicfunctions.CompletePaymentResponse{}, fmt.Errorf("subscription is not a string")
	}
	periodEndRaw, ok := input.Payload["current_period_end"].(float64)
	if !ok {
		return publicfunctions.CompletePaymentResponse{}, fmt.Errorf("current_period_end is not a float64")
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
		return publicfunctions.CompletePaymentResponse{}, err
	}
	if payment == nil {
		return publicfunctions.CompletePaymentResponse{}, myerror.NewPaymentNotFoundError()
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
	err = uc.txManager.Processing(ctx, func(ctx context.Context) error {
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
	if err != nil {
		return publicfunctions.CompletePaymentResponse{}, err
	}
	return publicfunctions.CompletePaymentResponse{
		PaymentID: payment.ID(),
		UserID:    payment.UserID(),
		Plan:      valueobject.NewSubscriptionPlanFromCode(valueobject.SubscriptionPlanPro),
		ExpiresAt: currentPeriodEnd,
	}, nil
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
