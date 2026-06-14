package purchasedomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)
type Subscription struct {
	id                  valueobject.SubscriptionID
	userID              valueobject.UserID
	plan                valueobject.SubscriptionPlan
	status              valueobject.SubscriptionStatus
	subscriptionOrderID *valueobject.SubscriptionOrderID // 手動作成はnullになる
	paymentID           valueobject.PaymentID
	activatedAt         time.Time
	expiresAt           time.Time
	canceledAt          *time.Time
	createdAt           time.Time
	updatedAt           time.Time
}

func (s *Subscription) ID() valueobject.SubscriptionID {
	return s.id
}

func (s *Subscription) UserID() valueobject.UserID {
	return s.userID
}

func (s *Subscription) Plan() valueobject.SubscriptionPlan {
	return s.plan
}

func (s *Subscription) Status() valueobject.SubscriptionStatus {
	return s.status
}

func (s *Subscription) SubscriptionOrderID() *valueobject.SubscriptionOrderID {
	return s.subscriptionOrderID
}

func (s *Subscription) PaymentID() valueobject.PaymentID {
	return s.paymentID
}

func (s *Subscription) ActivatedAt() time.Time {
	return s.activatedAt
}

func (s *Subscription) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Subscription) CanceledAt() *time.Time {
	return s.canceledAt
}

func (s *Subscription) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Subscription) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Subscription) MarkCanceled(at time.Time) {
	s.status = valueobject.NewSubscriptionStatusFromCode(valueobject.SubscriptionStatusCanceled)
	s.canceledAt = &at
	s.updatedAt = at
}

// ask : 期間経過時は何も更新しないで勝手に期限を過ぎる想定
func (s *Subscription) MarkExpired() {
	s.status = valueobject.NewSubscriptionStatusFromCode(valueobject.SubscriptionStatusExpired)
	s.updatedAt = time.Now()
}

func (s *Subscription) Renew(newExpiresAt time.Time) {
	s.expiresAt = newExpiresAt
	s.updatedAt = time.Now()
}

func CreateSubscription(userID valueobject.UserID, plan valueobject.SubscriptionPlan, orderID *valueobject.SubscriptionOrderID, paymentID valueobject.PaymentID, expiresAt time.Time) *Subscription {
	now := time.Now()
	return &Subscription{
		id:                  valueobject.NewPrimaryID[valueobject.SubscriptionID](),
		userID:              userID,
		plan:                plan,
		status:              valueobject.NewSubscriptionStatusFromCode(valueobject.SubscriptionStatusActive),
		subscriptionOrderID: orderID,
		paymentID:           paymentID,
		activatedAt:         now,
		expiresAt:           expiresAt,
		canceledAt:          nil,
		createdAt:           now,
		updatedAt:           now,
	}
}

func NewSubscription(id valueobject.SubscriptionID, userID valueobject.UserID, plan valueobject.SubscriptionPlan, status valueobject.SubscriptionStatus, subscriptionOrderID *valueobject.SubscriptionOrderID, paymentID valueobject.PaymentID, activatedAt time.Time, expiresAt time.Time, canceledAt *time.Time, createdAt time.Time, updatedAt time.Time) *Subscription {
	return &Subscription{
		id:                  id,
		userID:              userID,
		plan:                plan,
		status:              status,
		subscriptionOrderID: subscriptionOrderID,
		paymentID:           paymentID,
		activatedAt:         activatedAt,
		expiresAt:           expiresAt,
		canceledAt:          canceledAt,
		createdAt:           createdAt,
		updatedAt:           updatedAt,
	}
}
