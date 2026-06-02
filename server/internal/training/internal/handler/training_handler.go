package traininghandler

import (
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

type TrainingHandler struct {
	find   *trainingusecase.FindTrainingByID
	list   *trainingusecase.ListTrainings
	record *trainingusecase.RecordTraining
	update *trainingusecase.UpdateTraining
	delete *trainingusecase.DeleteTrainingByID
}

func New(
	find *trainingusecase.FindTrainingByID,
	list *trainingusecase.ListTrainings,
	record *trainingusecase.RecordTraining,
	update *trainingusecase.UpdateTraining,
	delete *trainingusecase.DeleteTrainingByID,
) http.Handler {
	h := &TrainingHandler{
		find:   find,
		list:   list,
		record: record,
		update: update,
		delete: delete,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /trainings/{id}", h.Find)
	mux.HandleFunc("GET /trainings", h.List)
	mux.HandleFunc("POST /trainings", h.Record)
	mux.HandleFunc("PUT /trainings/{id}", h.Update)
	mux.HandleFunc("DELETE /trainings/{id}", h.Delete)
	return mux
}

func (h *TrainingHandler) Find(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	trainingID, err := valueobject.NewPrimaryIDFromString[valueobject.TrainingID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid trainingID"))
		return
	}
	output, err := h.find.Execute(r.Context(), trainingusecase.FindTrainingByIDInput{UserID: userID, TrainingID: *trainingID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, trainingdto.NewTrainingDTO(output.Training))
}

func (h *TrainingHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	limit, offset := httpx.ParseOffsetPagination(r)
	output, err := h.list.Execute(r.Context(), trainingusecase.ListTrainingsInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	resp := trainingdto.ListTrainingsResponse{
		Trainings: lo.Map(output.Trainings, func(t *trainingdomain.Training, _ int) trainingdto.TrainingDTO {
			return trainingdto.NewTrainingDTO(t)
		}),
		Pagination: shareddto.PaginationDTO(output.Pagination),
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *TrainingHandler) Record(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req trainingdto.RecordTrainingRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}

	spec, err := req.ToSpec()
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	output, err := h.record.Execute(r.Context(), trainingusecase.RecordTrainingInput{UserID: userID, TrainingSpec: spec})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, trainingdto.RecordTrainingResponse{TrainingID: output.TrainingID.Value()})
}

func (h *TrainingHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
}

func (h *TrainingHandler) Delete(w http.ResponseWriter, r *http.Request) {}
