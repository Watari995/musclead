// Package paymentusecase は payment 集約の business logic を提供する。
//
// usecase 一覧 (Phase 1 で実装):
//   - ParseWebhookEvent (薄い wrapper、 stripeClient.ParseWebhookEvent を呼ぶだけ)
//   - InitiatePayment (purchase から呼ばれる、 Stripe Checkout Session 作成)
//   - CompletePayment (Webhook 内、 TX で payments / payment_events / outbox を atomic 更新)
//   - CancelPayment (Webhook 内、 解約反映)
//   - RenewPayment (将来用、 月次更新 Webhook 受信時)
//   - HandleFailure (薄い wrapper、 RetryStrategy.OnFailure を呼ぶだけ)
//   - GetPayment (Query 系)
//
// 既存参考: internal/weight/internal/usecase/record_weight.go
package paymentusecase

import "errors"

// errNotImplemented は skeleton 状態で usecase が呼ばれた時に返す sentinel error。
// 各 usecase の実装が完了したら、 該当の return から削除する (Phase 1 完了時には全部消える想定)。
var errNotImplemented = errors.New("payment usecase not implemented")
