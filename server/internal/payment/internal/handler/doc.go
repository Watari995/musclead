// Package paymenthandler は payment module の HTTP handler を提供する。
//
// 主な responsibility:
//   - POST /payment/webhook  (Stripe Webhook 受信、 Phase 4 で実装)
//
// 既存参考: internal/weight/internal/handler/weight_handler.go
package paymenthandler
