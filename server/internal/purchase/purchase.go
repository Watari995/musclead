// Package purchase is the public facade of the purchase module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
//
// 設計 (ADR 0013): purchase / payment 分離
// 依存方向: purchase → payment (purchase が payment.Command を呼ぶ)
//
// purchase 集約:
//   - subscription_orders: 申込トリガー (1 回限り、 履歴)
//   - subscriptions: 継続的な権利状態 (active / canceled / expired)
package purchase

import (
	"net/http"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasehandler "github.com/Watari995/musclead/internal/purchase/internal/handler"
	purchaseinfra "github.com/Watari995/musclead/internal/purchase/internal/infra"
	purchaseusecase "github.com/Watari995/musclead/internal/purchase/internal/usecase"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	Handler http.Handler
}

func NewModule(dbmap *gorp.DbMap, paymentCommand publicfunctions.PaymentCommand, userQuery userpublicfunctions.UserQuery, priceIDByPlan map[valueobject.SubscriptionPlanCode]string) *Module {
	// repositoryを作成
	dbmap.AddTableWithName(purchaseinfra.SubscriptionOrderModel{}, "subscription_orders").SetKeys(false, "ID")
	dbmap.AddTableWithName(purchaseinfra.SubscriptionModel{}, "subscriptions").SetKeys(false, "ID")
	orderRepo := purchaseinfra.NewSubscriptionOrderRepository(dbmap)

	subscribe := purchaseusecase.NewSubscribe(orderRepo, paymentCommand, userQuery, priceIDByPlan)
	handler := purchasehandler.NewPurchaseHandler(subscribe)
	return &Module{
		Handler: handler,
	}
}

//
// 必要な引数:
//   - dbmap *gorp.DbMap
//   - paymentCommand publicfunctions (payment) ← Phase 1 で公開済み、 paymentModule.Command() を注入
//   - cfg Config{StripeProPriceID 等の env}
//
// 流れ:
//   1. gorp model 登録 (SubscriptionOrderModel / SubscriptionModel)
//   2. Repository × 2 (SubscriptionOrder / Subscription)
//   3. usecase (SubscribeToPro、 paymentCommand を注入)
//   4. handler + mux 登録
//   5. publicfunctions の PurchaseCommand 実装 (Phase 2 では Phase 後半 or Phase 3 で公開)
//
// 参考: internal/payment/payment.go
