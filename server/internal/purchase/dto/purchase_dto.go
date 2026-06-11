// Package dto は purchase module の HTTP 入出力 DTO を定義する。
package dto

// TODO: HTTP リクエスト / レスポンス DTO 定義 (User 実装)
//
// 例 (Pro 申込 API):
//   type SubscribeToProRequest struct {
//       // payment_id は body で受け取らない (内部生成)
//   }
//   type SubscribeToProResponse struct {
//       CheckoutURL string `json:"checkout_url"`   // クライアントを Stripe Checkout にリダイレクトさせる先
//   }
//
// 参考: internal/user/dto/user_dto.go
