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

type ExerciseHandler struct {
	find         *trainingusecase.FindExerciseByID
	findBestSets *trainingusecase.FindBestSetsByExerciseIDs
	list         *trainingusecase.ListExercises
	create       *trainingusecase.CreateExercise
	update       *trainingusecase.UpdateExercise
	delete       *trainingusecase.DeleteExerciseByID
	reorder      *trainingusecase.ReorderExercises
}

func NewExerciseHandler(
	find *trainingusecase.FindExerciseByID,
	findBestSets *trainingusecase.FindBestSetsByExerciseIDs,
	list *trainingusecase.ListExercises,
	create *trainingusecase.CreateExercise,
	update *trainingusecase.UpdateExercise,
	delete *trainingusecase.DeleteExerciseByID,
	reorder *trainingusecase.ReorderExercises,
) http.Handler {
	h := &ExerciseHandler{
		find:         find,
		findBestSets: findBestSets,
		list:         list,
		create:       create,
		update:       update,
		delete:       delete,
		reorder:      reorder,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /exercises/{id}", h.Find)
	mux.HandleFunc("GET /exercises/best-sets", h.FindBestSets)
	mux.HandleFunc("GET /exercises", h.List)
	mux.HandleFunc("POST /exercises", h.Create)
	mux.HandleFunc("POST /exercises/reorder", h.Reorder)
	mux.HandleFunc("PUT /exercises/{id}", h.Update)
	mux.HandleFunc("DELETE /exercises/{id}", h.Delete)
	return mux
}

// Find godoc
//
// @Summary エクササイズ取得
// @Tags exercises
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 ExerciseID"
// @Success 200 {object} trainingdto.ExerciseDTO
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /exercises/{id} [get]
func (h *ExerciseHandler) Find(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid exerciseID"))
		return
	}
	output, err := h.find.Execute(r.Context(), trainingusecase.FindExerciseByIDInput{UserID: userID, ID: *exerciseID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.ExerciseFromEntity(output.Exercise))
}

// FindBestSets godoc
//
// @Summary 種目の最高記録(最重量セット)一覧を取得
// @Tags exercises
// @Produce json
// @Security BearerAuth
// @Param exercise_ids query []string true "対象 ExerciseID 一覧" collectionFormat(multi)
// @Success 200 {object} trainingdto.ListBestSetsResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /exercises/best-sets [get]
func (h *ExerciseHandler) FindBestSets(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	exerciseIDs := make([]valueobject.ExerciseID, 0, len(r.URL.Query()["exercise_ids"]))
	for _, raw := range r.URL.Query()["exercise_ids"] {
		id, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](raw)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid exerciseID"))
			return
		}
		exerciseIDs = append(exerciseIDs, *id)
	}
	output, err := h.findBestSets.Execute(r.Context(), trainingusecase.FindBestSetsByExerciseIDsInput{UserID: userID, ExerciseIDs: exerciseIDs})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.ListBestSetsResponse{
		BestSets: lo.Map(output.BestSets, func(b *trainingdomain.BestSetView, _ int) trainingdto.BestSetDTO {
			return trainingdto.BestSetFromData(b)
		}),
	})
}

// List godoc
//
// @Summary エクササイズ一覧
// @Tags exercises
// @Produce json
// @Security BearerAuth
// @Param limit query int false "1ページの件数 (default: 20, max: 100)"
// @Param offset query int false "開始位置 (default: 0)"
// @Success 200 {object} trainingdto.ListExercisesResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /exercises [get]
func (h *ExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	limit, offset := httpx.ParseOffsetPagination(r)
	output, err := h.list.Execute(r.Context(), trainingusecase.ListExercisesInput{UserID: userID, Limit: limit, Offset: offset})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := trainingdto.ListExercisesResponse{
		Exercises: lo.Map(output.Exercises, func(e *trainingdomain.Exercise, _ int) trainingdto.ExerciseDTO {
			return trainingdto.ExerciseFromEntity(e)
		}),
		Pagination: shareddto.PaginationDTO(output.Pagination),
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Create godoc
//
// @Summary エクササイズ作成
// @Tags exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body trainingdto.UpsertExerciseRequest true "エクササイズ作成"
// @Success 201 {object} trainingdto.UpsertExerciseResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /exercises [post]
func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req trainingdto.UpsertExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	name, err := valueobject.NewString50(req.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	output, err := h.create.Execute(r.Context(), trainingusecase.CreateExerciseInput{UserID: userID, Name: *name})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, trainingdto.UpsertExerciseResponse{ID: output.ID.Value()})
}

// Reorder godoc
//
// @Summary エクササイズ並び替え
// @Description 種目マスタ全件を、 渡された exercise_ids の順序どおりに並び替える。
// @Tags exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body trainingdto.ReorderExercisesRequest true "並び替え後の ExerciseID 一覧(全件)"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /exercises/reorder [post]
func (h *ExerciseHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req trainingdto.ReorderExercisesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	orderedIDs := make([]valueobject.ExerciseID, 0, len(req.ExerciseIDs))
	for _, raw := range req.ExerciseIDs {
		id, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](raw)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid exerciseID"))
			return
		}
		orderedIDs = append(orderedIDs, *id)
	}
	if err := h.reorder.Execute(r.Context(), trainingusecase.ReorderExercisesInput{UserID: userID, OrderedIDs: orderedIDs}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}

// Update godoc
//
// @Summary エクササイズ更新
// @Tags exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 ExerciseID"
// @Param request body trainingdto.UpsertExerciseRequest true "エクササイズ更新"
// @Success 200 {object} trainingdto.UpsertExerciseResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /exercises/{id} [put]
func (h *ExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid exerciseID"))
		return
	}
	var req trainingdto.UpsertExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	name, err := valueobject.NewString50(req.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	output, err := h.update.Execute(r.Context(), trainingusecase.UpdateExerciseInput{UserID: userID, ID: *exerciseID, Name: *name})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.UpsertExerciseResponse{ID: output.ID.Value()})
}

// Delete godoc
//
// @Summary エクササイズ削除
// @Tags exercises
// @Security BearerAuth
// @Produce json
// @Param id path string true "対象 ExerciseID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse "training_exercises から参照されている時"
// @Router /exercises/{id} [delete]
func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid exerciseID"))
		return
	}
	if err := h.delete.Execute(r.Context(), trainingusecase.DeleteExerciseByIDInput{UserID: userID, ID: *exerciseID}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}
