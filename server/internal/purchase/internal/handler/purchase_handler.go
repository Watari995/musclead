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
// 設計 (ADR 0013, 0017):
//   - POST /purchase/subscribe: 認証必須、 plan を受け取り、 Subscribe usecase を呼ぶ
//   - Stripe Checkout URL を返却、 クライアントは window.location でリダイレクト
type PurchaseHandler struct {
	subscribe     *purchaseusecase.Subscribe
	priceIDByPlan map[valueobject.SubscriptionPlanCode]string
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
		httpx.WriteError(w, err)
		return
	}

	var req dto.SubscribeRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}

	// TODO (User 実装): plan を VO 化
	//   plan, err := valueobject.NewSubscriptionPlanFromString(req.Plan)
	//   if err != nil { httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid plan").Wrap(err)); return }
	//
	// TODO (User 実装): priceID を解決
	//   priceID, ok := h.priceIDByPlan[valueobject.SubscriptionPlanCode(plan.Value())]
	//   if !ok { httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("unsupported plan")); return }
	//
	// TODO (User 実装 / 課題): user の Email を取得
	//   現状 user module の publicfunctions に Email を返す Query が無い (Authenticate のみ)。
	//   対応案:
	//     A. user/interface/publicfunctions に GetUserByID(userID) を追加 → main.go で paymentModule に注入
	//     B. JWT claim に email を入れて httpx.UserEmailFromContext で取り出す
	//     C. Subscribe usecase 内で user.Query を呼ぶ (purchase が user に依存、 ADR 0013 的に微妙)
	//   推奨は A (user module に閉じた Query 公開)。
	//
	// TODO (User 実装): Amount を解決
	//   amount, _ := valueobject.NewNonNegativeInt(480)
	//
	// TODO (User 実装): Subscribe.Execute 呼び出し
	//   output, err := h.subscribe.Execute(r.Context(), purchaseusecase.SubscribeInput{
	//       UserID:  userID,
	//       Email:   email,    // ← Email 取得が解決したら詰める
	//       Amount:  *amount,
	//       PriceID: priceID,
	//       Plan:    *plan,
	//   })
	//   if err != nil { httpx.WriteError(w, err); return }
	//
	//   httpx.WriteJSON(w, http.StatusOK, dto.SubscribeResponse{CheckoutURL: output.CheckoutURL.Value()})

	_ = userID
	httpx.WriteError(w, myerror.NewInternalError().SetMessage("purchase handler not implemented"))
}

func NewPurchaseHandler(subscribe *purchaseusecase.Subscribe, priceIDByPlan map[valueobject.SubscriptionPlanCode]string) *PurchaseHandler {
	return &PurchaseHandler{
		subscribe:     subscribe,
		priceIDByPlan: priceIDByPlan,
	}
}
