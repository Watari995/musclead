// Package purchaseusecase は purchase 集約の business logic を提供する。
//
// usecase 一覧 (Phase 2 で実装):
//   - SubscribeToPro: ⭐ 核心 orchestrator
//       1. subscription_orders INSERT (pending)
//       2. paymentCommand.InitiatePayment(...) ← Phase 1 で公開済み
//       3. subscription_orders UPDATE (payment_id をセット)
//       4. return CheckoutURL
//
// 将来 (Phase 9 outbox 受信 worker から呼ばれる):
//   - ActivateSubscription: Webhook 'checkout.session.completed' 後、 subscriptions INSERT
//   - CancelSubscription:   Webhook 'customer.subscription.updated' (cancel_at_period_end=true) 後
//   - ExpireSubscription:   Webhook 'customer.subscription.deleted' 後
//
// 参考: internal/payment/internal/usecase/initiate_payment.go (publicfunctions 直接実装パターン)
package purchaseusecase
