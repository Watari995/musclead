package paymentdomain

import (
	"context"
	"errors"

	"github.com/Watari995/musclead/internal/valueobject"
)

var ErrPaymentNotFound = errors.New("payment not found")

// PaymentRepository は payments テーブルへの永続化を抽象化する。
//
// 設計:
//   - 見つからない場合は (nil, nil) を返す (musclead 既存流儀、 sentinel error は使わない)
//   - Save は INSERT / UPDATE 兼用 (upsert) として実装する (既存 weight 流儀)
type PaymentRepository interface {
	// FindByID は payment_id (内部 UUID) で 1 件取得する。
	FindByID(ctx context.Context, id valueobject.PaymentID) (*Payment, error)

	// FindLatestSucceededByUserID は user の最新の succeeded payment を 1 件返す。
	// ADR 0017: 既存の Stripe Customer ID を再利用して二重作成を防ぐために使う。
	FindLatestSucceededByUserID(ctx context.Context, userID valueobject.UserID) (*Payment, error)

	// FindByStripeSubscriptionID は Stripe Subscription ID (sub_xxx) で payment を引く。
	// ADR 0014: 月次更新 / 解約 Webhook で「どの payment か」 を特定するために使う。
	FindByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*Payment, error)

	// Save は payment を保存する (INSERT / UPDATE 兼用)。
	Save(ctx context.Context, payment *Payment) error
}
