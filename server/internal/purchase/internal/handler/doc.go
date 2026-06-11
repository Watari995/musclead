// Package purchasehandler は purchase module の HTTP handler を提供する。
//
// 主な responsibility (Phase 2 で実装):
//   - POST /purchase/subscribe  (Pro 申込み、 認証必須)
//
// 将来:
//   - POST /purchase/portal-session (Customer Portal リダイレクト)
//
// 参考: internal/payment/internal/handler/webhook_handler.go (musclead 流儀の handler)
package purchasehandler
