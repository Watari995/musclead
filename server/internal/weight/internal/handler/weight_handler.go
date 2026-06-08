package weighthandler

import (
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdto "github.com/Watari995/musclead/internal/weight/dto"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
	weightusecase "github.com/Watari995/musclead/internal/weight/internal/usecase"
	"github.com/samber/lo"
)

type WeightHandler struct {
	record *weightusecase.RecordWeight
	find   *weightusecase.FindWeightByID
	list   *weightusecase.ListWeights
	update *weightusecase.UpdateWeight
	delete *weightusecase.DeleteWeightByID
}

func New(
	record *weightusecase.RecordWeight,
	find *weightusecase.FindWeightByID,
	list *weightusecase.ListWeights,
	update *weightusecase.UpdateWeight,
	delete *weightusecase.DeleteWeightByID,
) http.Handler {
	h := &WeightHandler{
		record: record,
		find:   find,
		list:   list,
		update: update,
		delete: delete,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /weights", h.Record)
	mux.HandleFunc("GET /weights/{id}", h.Find)
	mux.HandleFunc("GET /weights", h.List)
	mux.HandleFunc("PUT /weights/{id}", h.Update)
	mux.HandleFunc("DELETE /weights/{id}", h.Delete)
	return mux
}

// Record godoc
//
// @Summary 体重記録
// @Tags weights
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body weightdto.RecordWeightRequest true "体重記録"
// @Success 201 {object} weightdto.RecordWeightResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /weights [post]
func (h *WeightHandler) Record(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req weightdto.UpsertWeightRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	spec, err := req.ToSpec()
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	output, err := h.record.Execute(r.Context(), weightusecase.RecordWeightInput{
		UserID:     userID,
		WeightSpec: spec,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, weightdto.UpsertWeightResponse{WeightID: output.WeightID.Value()})
}

// Find godoc
//
// @Summary 体重取得
// @Tags weights
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 WeightID"
// @Success 200 {object} weightdto.WeightDTO
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /weights/{id} [get]
func (h *WeightHandler) Find(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	weightID, err := valueobject.NewPrimaryIDFromString[valueobject.WeightID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid weightID"))
		return
	}
	output, err := h.find.Execute(r.Context(), weightusecase.FindWeightByIDInput{
		ID:     *weightID,
		UserID: userID,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, weightdto.FromEntity(output.Weight))
}

// List godoc
//
// @Summary 体重一覧
// @Tags weights
// @Produce json
// @Security BearerAuth
// @Param limit query int false "1ページの件数 (default: 20, max: 100)"
// @Param offset query int false "開始位置 (default: 0)"
// @Success 200 {object} weightdto.ListWeightsResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /weights [get]
func (h *WeightHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	limit, offset := httpx.ParseOffsetPagination(r)
	output, err := h.list.Execute(r.Context(), weightusecase.ListWeightsInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := weightdto.ListWeightsResponse{
		Weights: lo.Map(output.Weights, func(w *weightdomain.Weight, _ int) weightdto.WeightDTO {
			return weightdto.FromEntity(w)
		}),
		Pagination: shareddto.NewPaginationDTO(output.Pagination),
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
//
// @Summary 体重更新
// @Tags weights
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 WeightID"
// @Param request body weightdto.UpsertWeightRequest true "体重更新"
// @Success 200 {object} weightdto.UpsertWeightResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /weights/{id} [put]
func (h *WeightHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	weightID, err := valueobject.NewPrimaryIDFromString[valueobject.WeightID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid weightID"))
		return
	}
	var req weightdto.UpsertWeightRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	spec, err := req.ToSpec()
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	output, err := h.update.Execute(r.Context(), weightusecase.UpdateWeightInput{
		ID:         *weightID,
		UserID:     userID,
		WeightSpec: spec,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, weightdto.UpsertWeightResponse{WeightID: output.WeightID.Value()})
}

// Delete godoc
//
// @Summary 体重削除
// @Tags weights
// @Security BearerAuth
// @Param id path string true "対象 WeightID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /weights/{id} [delete]
func (h *WeightHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	weightID, err := valueobject.NewPrimaryIDFromString[valueobject.WeightID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid weightID"))
		return
	}
	if err := h.delete.Execute(r.Context(), weightusecase.DeleteWeightByIDInput{
		ID:     *weightID,
		UserID: userID,
	}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}
