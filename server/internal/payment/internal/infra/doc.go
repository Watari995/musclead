// Package paymentinfra は payment domain interface の実装 (gorp / Stripe SDK / 等) を提供する。
//
// 含むもの:
//   - gorp model + repository: payment_models.go, payment_repository.go 等
//   - StripeClient impl: stripe_client.go (stripe-go v8x ラッパー)
//   - RetryStrategy impl: external_retry_strategy.go
//
// 既存参考: internal/weight/internal/infra/weight_repository.go (gorp パターン)
package paymentinfra
