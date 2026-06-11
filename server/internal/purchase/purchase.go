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
	dbmap.AddTableWithName(purchaseinfra.SubscriptionOrderModel{}, "subscription_orders").SetKeys(false, "ID")
	dbmap.AddTableWithName(purchaseinfra.SubscriptionModel{}, "subscriptions").SetKeys(false, "ID")
	orderRepo := purchaseinfra.NewSubscriptionOrderRepository(dbmap)

	subscribe := purchaseusecase.NewSubscribe(orderRepo, paymentCommand, userQuery, priceIDByPlan)
	handler := purchasehandler.NewPurchaseHandler(subscribe)
	return &Module{
		Handler: handler,
	}
}
