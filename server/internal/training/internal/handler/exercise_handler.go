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
	find   *trainingusecase.FindExerciseByID
	list   *trainingusecase.ListExercises
	create *trainingusecase.CreateExercise
	update *trainingusecase.UpdateExercise
	delete *trainingusecase.DeleteExerciseByID
}

func NewExerciseHandler(
	find *trainingusecase.FindExerciseByID,
	list *trainingusecase.ListExercises,
	create *trainingusecase.CreateExercise,
	update *trainingusecase.UpdateExercise,
	delete *trainingusecase.DeleteExerciseByID,
) http.Handler {
	h := &ExerciseHandler{
		find:   find,
		list:   list,
		create: create,
		update: update,
		delete: delete,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /exercises/{id}", h.Find)
	mux.HandleFunc("GET /exercises", h.List)
	mux.HandleFunc("POST /exercises", h.Create)
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
// @Success 200 {object} trainingdto.TrainingDTO
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
	httpx.WriteJSON(w, http.StatusOK, trainingdto.NewExerciseDTO(output.Exercise))
}

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
			return trainingdto.NewExerciseDTO(e)
		}),
		Pagination: shareddto.PaginationDTO(output.Pagination),
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

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
