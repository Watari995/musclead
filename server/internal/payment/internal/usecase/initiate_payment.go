package paymentusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// InitiatePayment は Pro 申込開始 (Stripe Checkout Session 作成) を担う。
//
// 流れ:
//  1. 既存 succeeded payment を user で検索 → Stripe Customer ID を再利用 (ADR 0017)
//  2. 新規 user なら StripeClient.CreateCustomer で新規 Customer 作成
//  3. payments INSERT (pending) + payment_events INSERT ('initiated')
//  4. StripeClient.CreateCheckoutSession で URL 取得 (PaymentID を Idempotency-Key)
//  5. payments UPDATE (checkout_url + stripe_checkout_session_id)
//  6. CheckoutURL を返す
//
// publicfunctions.PaymentCommand interface を直接実装する (musclead 既存流儀)。
type InitiatePayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeClient     paymentdomain.StripeClient
}

func (uc *InitiatePayment) InitiatePayment(ctx context.Context, req publicfunctions.InitiatePaymentRequest) (publicfunctions.InitiatePaymentResponse, error) {
	existing, _ := uc.paymentRepo.FindLatestSucceededByUserID(ctx, req.UserID)
	var stripeCustomerID string
	if existing != nil && existing.StripeCustomerID() != nil {
		stripeCustomerID = *existing.StripeCustomerID()
	} else {
		var err error
		stripeCustomerID, err = uc.stripeClient.CreateCustomer(ctx, paymentdomain.CreateCustomerInput{
			UserID: req.UserID, Email: req.Email,
		})
		if err != nil {
			return publicfunctions.InitiatePaymentResponse{}, err
		}
	}
	payment := paymentdomain.CreatePayment(
		req.UserID, valueobject.NewCurrencyFromCode(valueobject.CurrencyJPY), &stripeCustomerID, nil, nil, nil,
	)
	if err := uc.paymentRepo.Save(ctx, payment); err != nil {
		return publicfunctions.InitiatePaymentResponse{}, err
	}
	metadata := valueobject.Metadata{
		"currency":           valueobject.CurrencyJPY,
		"stripe_customer_id": stripeCustomerID,
	}
	paymentEvent := paymentdomain.CreatePaymentEvent(
		payment.ID(), valueobject.NewPaymentEventTypeFromCode(valueobject.PaymentEventTypeInitiated), metadata,
	)
	if err := uc.paymentEventRepo.Create(ctx, paymentEvent); err != nil {
		return publicfunctions.InitiatePaymentResponse{}, err
	}
	sess, err := uc.stripeClient.CreateCheckoutSession(ctx, paymentdomain.CreateCheckoutSessionInput{
		CustomerID: stripeCustomerID, PriceID: req.PriceID, PaymentID: payment.ID(),
	})
	if err != nil {
		return publicfunctions.InitiatePaymentResponse{}, err
	}
	payment.SetCheckoutSession(sess.SessionID, sess.CheckoutSessionURL)
	if err := uc.paymentRepo.Save(ctx, payment); err != nil {
		return publicfunctions.InitiatePaymentResponse{}, err
	}
	return publicfunctions.InitiatePaymentResponse{
		PaymentID:   payment.ID(),
		CheckoutURL: sess.CheckoutSessionURL,
	}, nil
}

func NewInitiatePayment(
	paymentRepo paymentdomain.PaymentRepository,
	paymentEventRepo paymentdomain.PaymentEventRepository,
	stripeClient paymentdomain.StripeClient,
) *InitiatePayment {
	return &InitiatePayment{
		paymentRepo:      paymentRepo,
		paymentEventRepo: paymentEventRepo,
		stripeClient:     stripeClient,
	}
}
