package paymentusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
)

// webhookCommand は 4 つの Webhook 系 usecase を束ねて
// publicfunctions.PaymentWebhookCommand interface を満たす facade。
//
// 設計 (ADR 0019):
//   - billing.NewModule に渡せる 1 つの依存に集約する役割
//   - embedding ではなく明示委譲にする (usecase の型名とメソッド名が同じため、
//     embedding するとフィールド名とメソッド名が衝突して interface を満たせない)
//   - 各 usecase は 1 struct 1 メソッドのまま (musclead 流儀を崩さない)
//
// なぜ module ファイル (payment.go) ではなく usecase 側の別ファイルに置くか:
//   - module ファイルは「組み立てるだけ」の薄い Composition Root に保ちたい。
//     委譲メソッドの実装ロジックを混ぜると太る。
//   - 束ね役 (委譲ロジック) は、 束ねる対象の usecase の隣に置くのが凝集度として自然。
//
// 注意: 単一メソッドの facade (例: PaymentCommand) はこの束ね役を作らず usecase を直接代入する。
// 委譲ロジックが無く上記の利点が出ないため。 メソッドが 2 つ以上になった時点で本パターンに切り替える。
type webhookCommand struct {
	completePayment *CompletePayment
	cancelPayment   *CancelPayment
	renewPayment    *RenewPayment
	handleFailure   *HandleFailure
}

func NewWebhookCommand(
	completePayment *CompletePayment,
	cancelPayment *CancelPayment,
	renewPayment *RenewPayment,
	handleFailure *HandleFailure,
) publicfunctions.PaymentWebhookCommand {
	return &webhookCommand{
		completePayment: completePayment,
		cancelPayment:   cancelPayment,
		renewPayment:    renewPayment,
		handleFailure:   handleFailure,
	}
}

func (w *webhookCommand) CompletePayment(ctx context.Context, req publicfunctions.CompletePaymentRequest) (publicfunctions.CompletePaymentResponse, error) {
	return w.completePayment.CompletePayment(ctx, req)
}

func (w *webhookCommand) CancelPayment(ctx context.Context, req publicfunctions.CancelPaymentRequest) (publicfunctions.CancelPaymentResponse, error) {
	return w.cancelPayment.CancelPayment(ctx, req)
}

func (w *webhookCommand) RenewPayment(ctx context.Context, req publicfunctions.RenewPaymentRequest) (publicfunctions.RenewPaymentResponse, error) {
	return w.renewPayment.RenewPayment(ctx, req)
}

func (w *webhookCommand) HandleFailure(ctx context.Context, req publicfunctions.HandleFailureRequest) error {
	return w.handleFailure.HandleFailure(ctx, req)
}
