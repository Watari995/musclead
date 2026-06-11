package purchasedomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// SubscriptionOrder は「Pro 申込みを試みた」 という 1 回限りのイベント。 失敗してもログ残す。
//
// 設計 (ADR 0013):
//   - 申込トリガー集約、 ライフサイクル: pending → succeeded / failed
//   - payment_id は申込開始時点では NULL、 payment.InitiatePayment が成功した直後に SET
//   - subscriptions (権利状態) とは別集約
//
// migration: sql/migrations/000020_create_subscription_orders.up.sql
//
//   - field 候補:
//       id           valueobject.SubscriptionOrderID (※ primary_id.go に追加必要)
//       userID       valueobject.UserID
//       plan         valueobject.SubscriptionPlan (既存 VO 再利用)
//       status       valueobject.SubscriptionOrderStatus (※ VO 化、 enum: pending/succeeded/failed)
//       paymentID    *valueobject.PaymentID (nullable、 申込開始時は nil)
//       succeededAt  *time.Time
//       failedAt     *time.Time
//       createdAt    time.Time
//       updatedAt    time.Time
//   - 状態遷移メソッド:
//       AttachPayment(paymentID PaymentID): payment.InitiatePayment 成功直後に呼ぶ
//       MarkSucceeded(at time.Time):  Webhook 経由で成功通知時
//       MarkFailed(at time.Time):     失敗通知時
//   - CreateSubscriptionOrder(userID, plan) → 新規 (pending、 paymentID は nil)
//   - NewSubscriptionOrder(全 field) → DB 復元用
//
// 参考: internal/payment/internal/domain/payment.go (state machine + Mark 系メソッドパターン)

type SubscriptionOrder struct {
	id          valueobject.SubscriptionOrderID
	userID      valueobject.UserID
	plan        valueobject.SubscriptionPlan
	status      valueobject.SubscriptionOrderStatus
	paymentID   *valueobject.PaymentID
	succeededAt *time.Time
	failedAt    *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func (s *SubscriptionOrder) ID() valueobject.SubscriptionOrderID {
	return s.id
}

func (s *SubscriptionOrder) UserID() valueobject.UserID {
	return s.userID
}

func (s *SubscriptionOrder) Plan() valueobject.SubscriptionPlan {
	return s.plan
}

func (s *SubscriptionOrder) Status() valueobject.SubscriptionOrderStatus {
	return s.status
}

func (s *SubscriptionOrder) PaymentID() *valueobject.PaymentID {
	return s.paymentID
}

func (s *SubscriptionOrder) SucceededAt() *time.Time {
	return s.succeededAt
}

func (s *SubscriptionOrder) FailedAt() *time.Time {
	return s.failedAt
}

func (s *SubscriptionOrder) CreatedAt() time.Time {
	return s.createdAt
}

func (s *SubscriptionOrder) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *SubscriptionOrder) AttachPayment(paymentID valueobject.PaymentID) {
	s.paymentID = &paymentID
	s.updatedAt = time.Now()
}

func (s *SubscriptionOrder) MarkSucceeded(at time.Time) {
	s.status = valueobject.NewSubscriptionOrderStatusFromCode(valueobject.SubscriptionOrderStatusSucceeded)
	s.succeededAt = &at
	s.updatedAt = at
}

func (s *SubscriptionOrder) MarkFailed(at time.Time) {
	s.status = valueobject.NewSubscriptionOrderStatusFromCode(valueobject.SubscriptionOrderStatusFailed)
	s.failedAt = &at
	s.updatedAt = at
}

func CreateSubscriptionOrder(userID valueobject.UserID, plan valueobject.SubscriptionPlan) *SubscriptionOrder {
	now := time.Now()
	return &SubscriptionOrder{
		id:          valueobject.NewPrimaryID[valueobject.SubscriptionOrderID](),
		userID:      userID,
		plan:        plan,
		status:      valueobject.NewSubscriptionOrderStatusFromCode(valueobject.SubscriptionOrderStatusPending),
		paymentID:   nil,
		succeededAt: nil,
		failedAt:    nil,
		createdAt:   now,
		updatedAt:   now,
	}
}

func NewSubscriptionOrder(id valueobject.SubscriptionOrderID, userID valueobject.UserID, plan valueobject.SubscriptionPlan, status valueobject.SubscriptionOrderStatus, paymentID *valueobject.PaymentID, succeededAt *time.Time, failedAt *time.Time, createdAt time.Time, updatedAt time.Time) *SubscriptionOrder {
	return &SubscriptionOrder{
		id:          id,
		userID:      userID,
		plan:        plan,
		status:      status,
		paymentID:   paymentID,
		succeededAt: succeededAt,
		failedAt:    failedAt,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
