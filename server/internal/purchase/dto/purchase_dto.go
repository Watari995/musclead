// Package dto は purchase module の HTTP 入出力 DTO を定義する。
package dto

// SubscribeRequest は POST /purchase/subscribe の body。
type SubscribeRequest struct {
	Plan string `json:"plan"` // "pro" 等、 valueobject.SubscriptionPlan で validate
}

// SubscribeResponse は申込み成功時のレスポンス。 client は CheckoutURL に遷移する。
type SubscribeResponse struct {
	CheckoutURL string `json:"checkout_url"`
}
