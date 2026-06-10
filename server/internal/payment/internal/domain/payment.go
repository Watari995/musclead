package paymentdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// Payment は payments テーブルに対応する集約。
// field 順は migration (sql/migrations/000015_create_payments.up.sql) の column 順に揃える。
type Payment struct {
	id                      valueobject.PaymentID
	userID                  valueobject.UserID
	amount                  valueobject.NonNegativeInt
	currency                valueobject.Currency
	status                  valueobject.PaymentStatus
	stripeCustomerID        *string
	stripeSubscriptionID    *string
	stripeCheckoutSessionID *string
	checkoutURL             *valueobject.URL
	currentPeriodEnd        *time.Time
	succeededAt             *time.Time
	failedAt                *time.Time
	failureReason           *string
	createdAt               time.Time
	updatedAt               time.Time
}

func (p *Payment) ID() valueobject.PaymentID {
	return p.id
}

func (p *Payment) UserID() valueobject.UserID {
	return p.userID
}

func (p *Payment) Amount() valueobject.NonNegativeInt {
	return p.amount
}

func (p *Payment) Currency() valueobject.Currency {
	return p.currency
}

func (p *Payment) Status() valueobject.PaymentStatus {
	return p.status
}

func (p *Payment) StripeCustomerID() *string {
	return p.stripeCustomerID
}

func (p *Payment) StripeSubscriptionID() *string {
	return p.stripeSubscriptionID
}

func (p *Payment) StripeCheckoutSessionID() *string {
	return p.stripeCheckoutSessionID
}

func (p *Payment) CheckoutURL() *valueobject.URL {
	return p.checkoutURL
}

func (p *Payment) CurrentPeriodEnd() *time.Time {
	return p.currentPeriodEnd
}

func (p *Payment) SucceededAt() *time.Time {
	return p.succeededAt
}

func (p *Payment) FailedAt() *time.Time {
	return p.failedAt
}

func (p *Payment) FailureReason() *string {
	return p.failureReason
}

func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Payment) UpdatedAt() time.Time {
	return p.updatedAt
}

func CreatePayment(
	userID valueobject.UserID,
	amount valueobject.NonNegativeInt,
	currency valueobject.Currency,
	stripeCustomerID *string,
	stripeCheckoutSessionID *string,
	checkoutURL *valueobject.URL,
	currentPeriodEnd *time.Time,
) *Payment {
	now := time.Now()
	return &Payment{
		id:                      valueobject.NewPrimaryID[valueobject.PaymentID](),
		userID:                  userID,
		amount:                  amount,
		currency:                currency,
		status:                  valueobject.NewPaymentStatusFromCode(valueobject.PaymentStatusPending),
		stripeCustomerID:        stripeCustomerID,
		stripeSubscriptionID:    nil,
		stripeCheckoutSessionID: stripeCheckoutSessionID,
		checkoutURL:             checkoutURL,
		currentPeriodEnd:        currentPeriodEnd,
		createdAt:               now,
		updatedAt:               now,
	}
}

func NewPayment(
	id valueobject.PaymentID,
	userID valueobject.UserID,
	amount valueobject.NonNegativeInt,
	currency valueobject.Currency,
	status valueobject.PaymentStatus,
	stripeCustomerID *string,
	stripeSubscriptionID *string,
	stripeCheckoutSessionID *string,
	checkoutURL *valueobject.URL,
	currentPeriodEnd *time.Time,
	succeededAt *time.Time,
	failedAt *time.Time,
	failureReason *string,
	createdAt time.Time,
	updatedAt time.Time,
) *Payment {
	return &Payment{
		id:                      id,
		userID:                  userID,
		amount:                  amount,
		currency:                currency,
		status:                  status,
		stripeCustomerID:        stripeCustomerID,
		stripeSubscriptionID:    stripeSubscriptionID,
		stripeCheckoutSessionID: stripeCheckoutSessionID,
		checkoutURL:             checkoutURL,
		currentPeriodEnd:        currentPeriodEnd,
		succeededAt:             succeededAt,
		failedAt:                failedAt,
		failureReason:           failureReason,
		createdAt:               createdAt,
		updatedAt:               updatedAt,
	}
}
