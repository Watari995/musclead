package purchasedomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// Subscription は「Pro である権利」 の継続的な状態。 Pro 機能 gate (将来) は本集約を見る。
//
// 設計 (ADR 0013, 0017):
//   - 権利状態集約、 ライフサイクル: active → canceled → expired
//   - Pro 判定: status に関係なく expires_at > NOW() なら Pro 扱い
//       - active + future:   Pro 利用可
//       - canceled + future: Pro 利用可 (解約予約中、 UI で「期末まで利用可」)
//       - expired or past expires_at: Pro 終了
//   - subscription_order_id は admin 手動作成の余地で nullable
//   - payment_id は Webhook で INSERT する時点で必ず存在するため NOT NULL
//
// migration: sql/migrations/000021_create_subscriptions.up.sql
//
// TODO (User 実装):
//   - field 候補:
//       id                    valueobject.SubscriptionID (※ primary_id.go に追加必要)
//       userID                valueobject.UserID
//       plan                  valueobject.SubscriptionPlan
//       status                valueobject.SubscriptionStatus (※ VO 化、 enum: active/canceled/expired)
//       subscriptionOrderID   *valueobject.SubscriptionOrderID  (nullable)
//       paymentID             valueobject.PaymentID             (NOT NULL)
//       activatedAt           time.Time
//       expiresAt             time.Time (NOT NULL)
//       canceledAt            *time.Time
//       createdAt             time.Time
//       updatedAt             time.Time
//   - 状態遷移メソッド:
//       MarkCanceled(at time.Time):   解約予約時 (Webhook customer.subscription.updated / Customer Portal)
//       MarkExpired():                期末経過時 (Webhook customer.subscription.deleted)
//       Renew(newExpiresAt time.Time): 月次更新時 (Webhook invoice.payment_succeeded)
//   - IsActive() bool: expires_at.After(time.Now()) で判定 (status は見ない、 ADR 0017)
//   - CreateSubscription(userID, plan, orderID, paymentID, expiresAt) → 新規 active
//   - NewSubscription(全 field) → DB 復元用
//
// 参考: internal/payment/internal/domain/payment.go

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
