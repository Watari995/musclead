package paymentusecase

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// InitiatePaymentInput は purchase 集約から呼ばれる入力。
type InitiatePaymentInput struct {
	UserID  valueobject.UserID
	Email   valueobject.Email          // Stripe Customer 作成時に渡す
	Amount  valueobject.NonNegativeInt // 480 (税込 JPY)
	PriceID string                     // Stripe Price ID (商品差分、 usecase 引数で受け取る)
}

// InitiatePaymentOutput は purchase 集約に返す結果。
type InitiatePaymentOutput struct {
	PaymentID   valueobject.PaymentID
	CheckoutURL valueobject.URL
}

// InitiatePayment は Pro 申込開始 (Stripe Checkout Session 作成) を担う。
//
// 流れ:
//  1. 既存 succeeded payment を user で検索 → Stripe Customer ID を再利用 (ADR 0017)
//  2. 新規 user なら StripeClient.CreateCustomer で新規 Customer 作成
//  3. payments INSERT (pending) + payment_events INSERT ('initiated')
//  4. StripeClient.CreateCheckoutSession で URL 取得 (PaymentID を Idempotency-Key)
//  5. payments UPDATE (checkout_url + stripe_checkout_session_id)
//  6. CheckoutURL を返す
type InitiatePayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeClient     paymentdomain.StripeClient
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

// Execute は申込を開始する。
//
// TODO (User 実装):
//
//	// 1. 既存 Stripe Customer ID 再利用
//	existing, _ := uc.paymentRepo.FindLatestSucceededByUserID(ctx, input.UserID)
//	var customerID string
//	if existing != nil && existing.StripeCustomerID() != nil {
//	    customerID = *existing.StripeCustomerID()
//	} else {
//	    cid, err := uc.stripeClient.CreateCustomer(ctx, paymentdomain.CreateCustomerInput{
//	        UserID: input.UserID, Email: input.Email,
//	    })
//	    if err != nil { return InitiatePaymentOutput{}, err }
//	    customerID = cid
//	}
//
//	// 2. payment INSERT (pending)
//	payment := paymentdomain.CreatePayment(input.UserID, input.Amount, currency, &customerID, nil, nil, nil)
//	if err := uc.paymentRepo.Save(ctx, payment); err != nil { ... }
//
//	// 3. Checkout Session 作成
//	sess, err := uc.stripeClient.CreateCheckoutSession(ctx, paymentdomain.CreateCheckoutSessionInput{
//	    CustomerID: customerID, PriceID: input.PriceID, PaymentID: payment.ID(),
//	})
//
//	// 4. payment UPDATE + checkout_url を返す
func (uc *InitiatePayment) Execute(ctx context.Context, input InitiatePaymentInput) (InitiatePaymentOutput, error) {
	existing, _ := uc.paymentRepo.FindLatestSucceededByUserID(ctx, input.UserID)
	var stripeCustomerID string
	if existing != nil && existing.StripeCustomerID() != nil {
		stripeCustomerID = *existing.StripeCustomerID()
	} else {
		var err error
		stripeCustomerID, err = uc.stripeClient.CreateCustomer(ctx, paymentdomain.CreateCustomerInput{
			UserID: input.UserID, Email: input.Email,
		})
		if err != nil {
			return InitiatePaymentOutput{}, err
		}
	}
	// payment INSERT (pending)
	payment := paymentdomain.CreatePayment(
		input.UserID, input.Amount, valueobject.NewCurrencyFromCode(valueobject.CurrencyJPY), &stripeCustomerID, nil, nil, nil,
	)
	if err := uc.paymentRepo.Save(ctx, payment); err != nil {
		return InitiatePaymentOutput{}, err
	}
	// payment_events INSERT (initiated)
	metadata := valueobject.Metadata{
		"amount":             input.Amount.Value(),
		"currency":           valueobject.CurrencyJPY,
		"stripe_customer_id": stripeCustomerID,
	}
	paymentEvent := paymentdomain.CreatePaymentEvent(
		payment.ID(), valueobject.NewPaymentEventTypeFromCode(valueobject.PaymentEventTypeInitiated), metadata,
	)
	if err := uc.paymentEventRepo.Create(ctx, paymentEvent); err != nil {
		return InitiatePaymentOutput{}, err
	}
	// Checkout Session Created
	sess, err := uc.stripeClient.CreateCheckoutSession(ctx, paymentdomain.CreateCheckoutSessionInput{
		CustomerID: stripeCustomerID, PriceID: input.PriceID, PaymentID: payment.ID(),
	})
	if err != nil {
		return InitiatePaymentOutput{}, err
	}
	// payment update
	payment.SetCheckoutSession(sess.SessionID, sess.CheckoutSessionURL)
	if err := uc.paymentRepo.Save(ctx, payment); err != nil {
		return InitiatePaymentOutput{}, err
	}
	return InitiatePaymentOutput{
		PaymentID:   payment.ID(),
		CheckoutURL: sess.CheckoutSessionURL,
	}, nil
}
