// Package publicfunctions は payment module が他 module に公開する Command / Query interface を定義する。
// 他の module (purchase 等) はこのパッケージにのみ依存する。
//
// 設計: ADR 0013 で確定の Facade パターン
//   - InitiatePayment    (購入オーケストレータから「決済開始してくれ」 と呼ばれる)
//   - CompletePayment    (Webhook handler から「決済完了 (succeeded / failed) を反映」 と呼ばれる)
//   - CancelPayment      (Customer Portal 経由の解約反映)
//   - CapturePayment     (サブスクでは未使用、 将来用)
package publicfunctions

import (
	"context"
)

// TODO: 各 Command の Request / Response struct を定義
// 参考: internal/user/interface/publicfunctions/command.go の AuthenticateRequest/Response パターン
//
// 例: InitiatePaymentRequest は (UserID, Amount, Currency, Plan 等) を含む

// PaymentCommand は payment 集約に対する書き込み系操作の公開 API。
// purchase context からは Module.Command() 経由でアクセスする。
type PaymentCommand interface {
	// TODO: InitiatePayment(ctx context.Context, req InitiatePaymentRequest) (InitiatePaymentResponse, error)
	// TODO: CompletePayment(ctx context.Context, req CompletePaymentRequest) error
	// TODO: CancelPayment(ctx context.Context, req CancelPaymentRequest) error
	// TODO: CapturePayment は将来用、 MVP では不要

	// Placeholder: interface 文法エラー回避のため、 後で TODO に置き換える
	_placeholder(ctx context.Context) error
}
