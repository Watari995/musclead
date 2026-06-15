package purchaseusecase

import (
	"context"

	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
)

// purchaseCommand は複数の usecase を束ねて publicfunctions.PurchaseCommand を満たす。
//
// 設計 (payment の webhookCommand と同じ Case C):
//   - 各 usecase は「型名 == メソッド名」 (ActivateSubscription.ActivateSubscription 等) なので
//     embed すると フィールド名 と 昇格メソッド名 が衝突する。
//   - よって embed せず、 名前付きフィールド + 明示的な委譲メソッドにする。
//
// なぜ module ファイル (purchase.go) ではなく usecase 側の別ファイルに置くか:
//   - module ファイルは「組み立てるだけ」の薄い Composition Root に保ちたい。
//     委譲メソッドの実装ロジックを混ぜると太る。
//   - 束ね役 (委譲ロジック) は、 束ねる対象の usecase の隣に置くのが凝集度として自然。
//
// 注意: 単一メソッドの facade はこの束ね役を作らず usecase を直接代入する (委譲ロジックが無く
// 上記の利点が出ないため)。 メソッドが 2 つ以上になった時点で本パターンに切り替える (ActivateSubscription
// だけだった頃は直接代入で、 RenewSubscription を足した時に本束ね役へ移行した)。
type purchaseCommand struct {
	activateSubscription *ActivateSubscription
	renewSubscription    *RenewSubscription
}

func NewPurchaseCommand(
	activateSubscription *ActivateSubscription,
	renewSubscription *RenewSubscription,
) purchasepublicfunctions.PurchaseCommand {
	return &purchaseCommand{
		activateSubscription: activateSubscription,
		renewSubscription:    renewSubscription,
	}
}

func (c *purchaseCommand) ActivateSubscription(ctx context.Context, req purchasepublicfunctions.ActivateSubscriptionRequest) error {
	return c.activateSubscription.ActivateSubscription(ctx, req)
}

func (c *purchaseCommand) RenewSubscription(ctx context.Context, req purchasepublicfunctions.RenewSubscriptionRequest) error {
	return c.renewSubscription.RenewSubscription(ctx, req)
}
