// Package publicfunctions ... (Command と同パッケージ)
package publicfunctions

import (
	"context"
)

// PaymentQuery は payment 集約に対する読み込み系操作の公開 API。
// 例: /users/me で「最新の successful payment」 を取得する用途。
type PaymentQuery interface {
	// TODO: GetPayment(ctx, paymentID) (PaymentDTO, error)
	// TODO: FindLatestByUserID(ctx, userID) (*PaymentDTO, error)
	//   - 既存 Stripe Customer ID を再利用するため (ADR 0017)
	//   - 未加入なら nil 返却

	// Placeholder
	_placeholderQuery(ctx context.Context) error
}
