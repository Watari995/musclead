package traininghandler

import (
	"encoding/json"
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/Watari995/musclead/internal/shared/httpx"
	trainingdto "github.com/Watari995/musclead/internal/training/dto"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
)

type RoutineHandler struct {
	find    *trainingusecase.FindRoutineByID
	list    *trainingusecase.ListRoutines
	create  *trainingusecase.CreateRoutine
	update  *trainingusecase.UpdateRoutine
	delete  *trainingusecase.DeleteRoutineByID
	reorder *trainingusecase.ReorderRoutines
}

func NewRoutineHandler(
	find *trainingusecase.FindRoutineByID,
	list *trainingusecase.ListRoutines,
	create *trainingusecase.CreateRoutine,
	update *trainingusecase.UpdateRoutine,
	delete *trainingusecase.DeleteRoutineByID,
	reorder *trainingusecase.ReorderRoutines,
) http.Handler {
	h := &RoutineHandler{
		find:    find,
		list:    list,
		create:  create,
		update:  update,
		delete:  delete,
		reorder: reorder,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /routines/{id}", h.Find)
	mux.HandleFunc("GET /routines", h.List)
	mux.HandleFunc("POST /routines", h.Create)
	mux.HandleFunc("POST /routines/reorder", h.Reorder)
	mux.HandleFunc("PUT /routines/{id}", h.Update)
	mux.HandleFunc("DELETE /routines/{id}", h.Delete)
	return mux
}

// Find godoc
//
// @Summary ルーティン詳細取得
// @Tags routines
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 RoutineID"
// @Success 200 {object} trainingdto.RoutineDTO
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /routines/{id} [get]
func (h *RoutineHandler) Find(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	routineID, err := valueobject.NewPrimaryIDFromString[valueobject.RoutineID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid routineID"))
		return
	}
	output, err := h.find.Execute(r.Context(), trainingusecase.FindRoutineByIDInput{UserID: userID, ID: *routineID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.RoutineFromEntity(output.Routine))
}

// List godoc
//
// @Summary ルーティン一覧
// @Tags routines
// @Produce json
// @Security BearerAuth
// @Param limit query int false "1ページの件数 (default: 20, max: 100)"
// @Param offset query int false "開始位置 (default: 0)"
// @Success 200 {object} trainingdto.ListRoutinesResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /routines [get]
func (h *RoutineHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	limit, offset := httpx.ParseOffsetPagination(r)
	output, err := h.list.Execute(r.Context(), trainingusecase.ListRoutinesInput{UserID: userID, Limit: limit, Offset: offset})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.ListRoutinesResponse{
		Routines: lo.Map(output.Routines, func(r *trainingdomain.RoutineView, _ int) trainingdto.RoutineDTO {
			return trainingdto.RoutineFromEntity(r)
		}),
		Pagination: shareddto.PaginationDTO(output.Pagination),
	})
}

// Create godoc
//
// @Summary ルーティン作成
// @Tags routines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body trainingdto.UpsertRoutineRequest true "ルーティン作成"
// @Success 201 {object} trainingdto.UpsertRoutineResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse "同 user 内で名前重複"
// @Router /routines [post]
func (h *RoutineHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	input := trainingdto.UpsertRoutineRequest{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	name, err := valueobject.NewString50(input.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	specs, err := input.ToSpec()
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	output, err := h.create.Execute(r.Context(), trainingusecase.CreateRoutineInput{UserID: userID, RoutineSpec: trainingdomain.RoutineSpec{Name: *name, Exercises: specs}})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, trainingdto.UpsertRoutineResponse{ID: output.ID.Value()})
}

// Reorder godoc
//
// @Summary ルーティン並び替え
// @Description ルーティン全件を、 渡された routine_ids の順序どおりに並び替える。
// @Tags routines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body trainingdto.ReorderRoutinesRequest true "並び替え後の RoutineID 一覧(全件)"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /routines/reorder [post]
func (h *RoutineHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req trainingdto.ReorderRoutinesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	orderedIDs := make([]valueobject.RoutineID, 0, len(req.RoutineIDs))
	for _, raw := range req.RoutineIDs {
		id, err := valueobject.NewPrimaryIDFromString[valueobject.RoutineID](raw)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid routineID"))
			return
		}
		orderedIDs = append(orderedIDs, *id)
	}
	if err := h.reorder.Execute(r.Context(), trainingusecase.ReorderRoutinesInput{UserID: userID, OrderedIDs: orderedIDs}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}

// Update godoc
//
// @Summary ルーティン更新
// @Tags routines
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 RoutineID"
// @Param request body trainingdto.UpsertRoutineRequest true "ルーティン更新"
// @Success 200 {object} trainingdto.UpsertRoutineResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse "同 user 内で名前重複"
// @Router /routines/{id} [put]
func (h *RoutineHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	routineID, err := valueobject.NewPrimaryIDFromString[valueobject.RoutineID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid routineID"))
		return
	}
	input := trainingdto.UpsertRoutineRequest{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	name, err := valueobject.NewString50(input.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	specs, err := input.ToSpec()
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	output, err := h.update.Execute(r.Context(), trainingusecase.UpdateRoutineInput{ID: *routineID, UserID: userID, RoutineSpec: trainingdomain.RoutineSpec{Name: *name, Exercises: specs}})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.UpsertRoutineResponse{ID: output.ID.Value()})
}

// Delete godoc
//
// @Summary ルーティン削除
// @Tags routines
// @Security BearerAuth
// @Produce json
// @Param id path string true "対象 RoutineID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /routines/{id} [delete]
func (h *RoutineHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	routineID, err := valueobject.NewPrimaryIDFromString[valueobject.RoutineID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid routineID"))
		return
	}
	if err := h.delete.Execute(r.Context(), trainingusecase.DeleteRoutineByIDInput{ID: *routineID, UserID: userID}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}
