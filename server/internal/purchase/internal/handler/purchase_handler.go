package purchasehandler

import (
	"net/http"

	purchaseusecase "github.com/Watari995/musclead/internal/purchase/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/httpx"
)

// PurchaseHandler は purchase 関連の HTTP handler を提供する。
//
// 設計 (ADR 0013):
//   - handler は HTTP I/O のみ。 publicfunctions / business 設定値は持たない
//   - business logic は Subscribe usecase に集約 (priceID 解決 / email 取得 / 金額決定 等)
type PurchaseHandler struct {
	subscribe *purchaseusecase.Subscribe
}

// Subscribe godoc
//
// @Summary Pro 申込み (Stripe Checkout 起動)
// @Tags purchase
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SubscribeRequest true "申込み plan"
// @Success 200 {object} dto.SubscribeResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /purchase/subscribe [post]
//
// 流れ:
//  1. UserID を context から取得 (httpx.UserIDFromContext)
//  2. body を dto.SubscribeRequest にデコード (httpx.DecodeJSON)
//  3. req.Plan を valueobject.SubscriptionPlan に validate
//  4. h.subscribe.Execute(ctx, SubscribeInput{UserID, Plan}) を呼ぶ
//  5. CheckoutURL を dto.SubscribeResponse として httpx.WriteJSON で返却
func (h *PurchaseHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	// TODO (User 実装): 上記の「流れ」 を実装
	httpx.WriteError(w, nil)
}

func NewPurchaseHandler(subscribe *purchaseusecase.Subscribe) *PurchaseHandler {
	return &PurchaseHandler{subscribe: subscribe}
}
