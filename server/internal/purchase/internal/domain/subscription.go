package purchasedomain

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
