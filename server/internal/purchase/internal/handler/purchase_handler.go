package purchasehandler

import (
	"net/http"
	"time"

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
	subscribe           *purchaseusecase.Subscribe
	getSubscription     *purchaseusecase.GetSubscription
	createPortalSession *purchaseusecase.CreatePortalSession
}

func NewPurchaseHandler(subscribe *purchaseusecase.Subscribe, getSubscription *purchaseusecase.GetSubscription, createPortalSession *purchaseusecase.CreatePortalSession) http.Handler {
	h := &PurchaseHandler{subscribe: subscribe, getSubscription: getSubscription, createPortalSession: createPortalSession}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /purchase/subscribe", h.Subscribe)
	mux.HandleFunc("GET /purchase/subscription", h.GetSubscription)
	mux.HandleFunc("POST /purchase/portal-session", h.CreatePortalSession)
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

// GetSubscription godoc
//
// @Summary サブスクリプション状態取得
// @Tags purchase
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.GetSubscriptionResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /purchase/subscription [get]
func (h *PurchaseHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, myerror.NewUnauthorizedError())
		return
	}
	output, err := h.getSubscription.Execute(r.Context(), purchaseusecase.GetSubscriptionInput{UserID: userID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var plan *string
	if output.Plan != nil {
		s := output.Plan.String()
		plan = &s
	}
	var expiresAt *string
	if output.ExpiresAt != nil {
		s := output.ExpiresAt.Format(time.RFC3339)
		expiresAt = &s
	}
	httpx.WriteJSON(w, http.StatusOK, dto.GetSubscriptionResponse{IsPro: output.IsPro, Plan: plan, ExpiresAt: expiresAt})
}

// CreatePortalSession godoc
//
// @Summary Customer Portal セッション作成
// @Tags purchase
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CreatePortalSessionResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /purchase/portal-session [post]
func (h *PurchaseHandler) CreatePortalSession(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, myerror.NewUnauthorizedError())
		return
	}
	output, err := h.createPortalSession.Execute(r.Context(), purchaseusecase.CreatePortalSessionInput{UserID: userID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto.CreatePortalSessionResponse{PortalURL: output.PortalURL.String()})
}
