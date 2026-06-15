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

type RenewPayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	stripeClient     paymentdomain.StripeClient
	txManager        dbtx.TransactionManager
}

func (uc *RenewPayment) RenewPayment(ctx context.Context, input publicfunctions.RenewPaymentRequest) (publicfunctions.RenewPaymentResponse, error) {
	stripeSubscriptionID, ok := input.Payload["subscription"].(string)
	if !ok {
		return publicfunctions.RenewPaymentResponse{}, fmt.Errorf("subscription is not a string")
	}
	// invoice.paid の payload 直下に current_period_end は無いため、 Stripe から取得する (権威ある値)。
	currentPeriodEnd, err := uc.stripeClient.RetrieveSubscription(ctx, stripeSubscriptionID)
	if err != nil {
		return publicfunctions.RenewPaymentResponse{}, myerror.NewInternalError().Wrap(err)
	}
	stripeEventMetadata := valueobject.Metadata{
		"stripe_event_id":        input.StripeEventID,
		"stripe_subscription_id": stripeSubscriptionID,
		"current_period_end":     currentPeriodEnd.Format(time.RFC3339),
	}
	stripeEvent := paymentdomain.CreateStripeEvent(input.StripeEventID, input.EventType, stripeEventMetadata)
	payment, err := uc.paymentRepo.FindByStripeSubscriptionID(ctx, stripeSubscriptionID)
	if err != nil {
		return publicfunctions.RenewPaymentResponse{}, myerror.NewInternalError().Wrap(err)
	}
	if payment == nil {
		return publicfunctions.RenewPaymentResponse{}, myerror.NewPaymentNotFoundError()
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
	if err := uc.txManager.Processing(ctx, func(ctx context.Context) error {
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
	}); err != nil {
		return publicfunctions.RenewPaymentResponse{}, err
	}
	return publicfunctions.RenewPaymentResponse{PaymentID: payment.ID(), ExpiresAt: currentPeriodEnd}, nil
}

func NewRenewPayment(
	paymentRepo paymentdomain.PaymentRepository,
	paymentEventRepo paymentdomain.PaymentEventRepository,
	stripeEventRepo paymentdomain.StripeEventRepository,
	outboxEventRepo paymentdomain.OutboxEventRepository,
	stripeClient paymentdomain.StripeClient,
	txManager dbtx.TransactionManager,
) *RenewPayment {
	return &RenewPayment{
		paymentRepo:      paymentRepo,
		paymentEventRepo: paymentEventRepo,
		stripeEventRepo:  stripeEventRepo,
		outboxEventRepo:  outboxEventRepo,
		stripeClient:     stripeClient,
		txManager:        txManager,
	}
}
