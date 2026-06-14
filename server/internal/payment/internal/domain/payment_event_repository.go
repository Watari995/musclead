package paymentdomain

import (
	"context"
)

// PaymentEventRepository は payment_events テーブルへの永続化を抽象化する。
//
// 設計:
//   - append-only な監査ログ。 状態を書き換えない (Save = INSERT のみ)
//   - 監査画面 / 履歴照会の Find 系メソッドは MVP では不要なので追加しない
type PaymentEventRepository interface {
	Create(ctx context.Context, event *PaymentEvent) error
}
