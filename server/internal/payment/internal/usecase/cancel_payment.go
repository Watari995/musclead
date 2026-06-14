package paymentusecase

import (
	"context"
	"errors"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CancelPayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

func (uc *CancelPayment) CancelPayment(ctx context.Context, input publicfunctions.CancelPaymentRequest) error {
	stripeSubscriptionID, ok := input.Payload["subscription"].(string)
	if !ok {
		return myerror.NewInternalError().SetMessage("subscription is not a string")
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
