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

	paymentpublicfunctions "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"

	purchasehandler "github.com/Watari995/musclead/internal/purchase/internal/handler"
	purchaseinfra "github.com/Watari995/musclead/internal/purchase/internal/infra"
	purchaseusecase "github.com/Watari995/musclead/internal/purchase/internal/usecase"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	Handler         http.Handler
	purchaseCommand purchasepublicfunctions.PurchaseCommand
}

func NewModule(dbmap *gorp.DbMap, paymentCommand paymentpublicfunctions.PaymentCommand, userQuery userpublicfunctions.UserQuery, priceIDByPlan map[valueobject.SubscriptionPlanCode]string) *Module {
	dbmap.AddTableWithName(purchaseinfra.SubscriptionOrderModel{}, "subscription_orders").SetKeys(false, "ID")
	dbmap.AddTableWithName(purchaseinfra.SubscriptionModel{}, "subscriptions").SetKeys(false, "ID")
	orderRepo := purchaseinfra.NewSubscriptionOrderRepository(dbmap)
	subscriptionRepo := purchaseinfra.NewSubscriptionRepository(dbmap)

	subscribe := purchaseusecase.NewSubscribe(orderRepo, paymentCommand, userQuery, priceIDByPlan)
	activateSubscription := purchaseusecase.NewActivateSubscription(subscriptionRepo, orderRepo)
	handler := purchasehandler.NewPurchaseHandler(subscribe)
	return &Module{
		Handler:         handler,
		purchaseCommand: activateSubscription,
	}
}

// PurchaseCommand は他 module 公開用 getter。
// purchaseCommand は unexported にし setter を持たないことで、 NewModule 後の依存差し替えを防ぐ。
func (m *Module) PurchaseCommand() purchasepublicfunctions.PurchaseCommand {
	return m.purchaseCommand
}
