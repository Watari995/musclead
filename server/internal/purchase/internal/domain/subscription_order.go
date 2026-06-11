package purchasedomain

// SubscriptionOrder は「Pro 申込みを試みた」 という 1 回限りのイベント。 失敗してもログ残す。
//
// 設計 (ADR 0013):
//   - 申込トリガー集約、 ライフサイクル: pending → succeeded / failed
//   - payment_id は申込開始時点では NULL、 payment.InitiatePayment が成功した直後に SET
//   - subscriptions (権利状態) とは別集約
//
// migration: sql/migrations/000020_create_subscription_orders.up.sql
//
// TODO (User 実装):
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
