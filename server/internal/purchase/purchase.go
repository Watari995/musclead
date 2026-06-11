// Package purchase is the public facade of the purchase module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
//
// 設計 (ADR 0013): purchase / payment 分離
// 依存方向: purchase → payment (purchase が payment.Command を呼ぶ)
//
// purchase 集約:
//   - subscription_orders: 申込トリガー (1 回限り、 履歴)
//   - subscriptions: 継続的な権利状態 (active / canceled / expired)
package purchase

// TODO: NewModule の実装 (User が wire を組み立てる)
//
// 必要な引数:
//   - dbmap *gorp.DbMap
//   - paymentCommand publicfunctions (payment) ← Phase 1 で公開済み、 paymentModule.Command() を注入
//   - cfg Config{StripeProPriceID 等の env}
//
// 流れ:
//   1. gorp model 登録 (SubscriptionOrderModel / SubscriptionModel)
//   2. Repository × 2 (SubscriptionOrder / Subscription)
//   3. usecase (SubscribeToPro、 paymentCommand を注入)
//   4. handler + mux 登録
//   5. publicfunctions の PurchaseCommand 実装 (Phase 2 では Phase 後半 or Phase 3 で公開)
//
// 参考: internal/payment/payment.go
