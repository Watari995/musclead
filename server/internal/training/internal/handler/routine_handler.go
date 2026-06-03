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
	find   *trainingusecase.FindRoutineByID
	list   *trainingusecase.ListRoutines
	create *trainingusecase.CreateRoutine
	update *trainingusecase.UpdateRoutine
	delete *trainingusecase.DeleteRoutineByID
}

func NewRoutineHandler(
	find *trainingusecase.FindRoutineByID,
	list *trainingusecase.ListRoutines,
	create *trainingusecase.CreateRoutine,
	update *trainingusecase.UpdateRoutine,
	delete *trainingusecase.DeleteRoutineByID,
) http.Handler {
	h := &RoutineHandler{
		find:   find,
		list:   list,
		create: create,
		update: update,
		delete: delete,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /routines/{id}", h.Find)
	mux.HandleFunc("GET /routines", h.List)
	mux.HandleFunc("POST /routines", h.Create)
	mux.HandleFunc("PUT /routines/{id}", h.Update)
	mux.HandleFunc("DELETE /routines/{id}", h.Delete)
	return mux
}

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
	httpx.WriteJSON(w, http.StatusOK, trainingdto.NewRoutineDTO(output.Routine))
}

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
			return trainingdto.NewRoutineDTO(r)
		}),
		Pagination: shareddto.PaginationDTO(output.Pagination),
	})
}

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
	specs := make([]trainingdomain.RoutineExerciseSpec, 0, len(input.Exercises))
	for _, e := range input.Exercises {
		exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](e.ExerciseID)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid exercise_id"))
			return
		}
		displayOrder, err := valueobject.NewNonNegativeInt(e.DisplayOrder)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid display_order"))
			return
		}
		specs = append(specs, trainingdomain.RoutineExerciseSpec{ExerciseID: *exerciseID, DisplayOrder: *displayOrder})
	}
	output, err := h.create.Execute(r.Context(), trainingusecase.CreateRoutineInput{UserID: userID, RoutineSpec: trainingdomain.RoutineSpec{Name: *name, Exercises: specs}})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, trainingdto.UpsertRoutineResponse{ID: output.ID.Value()})
}

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
	specs := make([]trainingdomain.RoutineExerciseSpec, 0, len(input.Exercises))
	for _, e := range input.Exercises {
		exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](e.ExerciseID)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid exercise_id"))
			return
		}
		displayOrder, err := valueobject.NewNonNegativeInt(e.DisplayOrder)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid display_order"))
			return
		}
		specs = append(specs, trainingdomain.RoutineExerciseSpec{ExerciseID: *exerciseID, DisplayOrder: *displayOrder})
	}
	output, err := h.update.Execute(r.Context(), trainingusecase.UpdateRoutineInput{ID: *routineID, UserID: userID, RoutineSpec: trainingdomain.RoutineSpec{Name: *name, Exercises: specs}})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.UpsertRoutineResponse{ID: output.ID.Value()})
}

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
