package purchasehandler

import (
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/purchase/dto"
	purchaseusecase "github.com/Watari995/musclead/internal/purchase/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
)

// PurchaseHandler は purchase 関連の HTTP handler を提供する。
//
// 設計 (ADR 0013):
//   - handler は HTTP I/O のみ。 publicfunctions / business 設定値は持たない
//   - business logic は Subscribe usecase に集約 (priceID 解決 / email 取得 / 金額決定 等)
type PurchaseHandler struct {
	subscribe *purchaseusecase.Subscribe
}

func NewPurchaseHandler(subscribe *purchaseusecase.Subscribe) http.Handler {
	h := &PurchaseHandler{subscribe: subscribe}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /purchase/subscribe", h.Subscribe)
	return mux
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
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, myerror.NewUnauthorizedError())
		return
	}
	var req dto.SubscribeRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	plan, err := valueobject.NewSubscriptionPlanFromString(req.Plan)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid plan"))
		return
	}
	output, err := h.subscribe.Execute(r.Context(), purchaseusecase.SubscribeInput{UserID: userID, Plan: *plan})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto.SubscribeResponse{CheckoutURL: output.CheckoutURL.String()})
}
