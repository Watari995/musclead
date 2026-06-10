// Package paymentusecase は payment 集約の business logic を提供する。
//
// usecase 一覧 (Phase 1 で実装):
//   - InitiatePayment (purchase から呼ばれる、 Stripe Checkout Session 作成)
//   - CompletePayment (Webhook 内で呼ばれる、 payments / payment_events / outbox を atomic に更新)
//   - CancelPayment (Customer Portal 経由の解約反映)
//   - GetPayment / FindLatestByUserID (Query 系)
//
// 既存参考: internal/weight/internal/usecase/record_weight.go
package paymentusecase
