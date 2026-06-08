package weighthandler

import (
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdto "github.com/Watari995/musclead/internal/weight/dto"
	weightusecase "github.com/Watari995/musclead/internal/weight/internal/usecase"
)

type WeightHandler struct {
	record *weightusecase.RecordWeight
}

func New(record *weightusecase.RecordWeight) http.Handler {
	h := &WeightHandler{
		record: record,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /weights", h.Record)
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
	var req weightdto.RecordWeightRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	weightKg, err := valueobject.NewWeightKgFromString(req.WeightKg)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid weight kg"))
		return
	}
	var bodyFatPercentage *valueobject.Percentage
	if req.BodyFatPercentage != nil {
		bodyFatPercentage, err = valueobject.NewPercentageFromString(*req.BodyFatPercentage)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid body fat percentage"))
			return
		}
	}
	var skeletalMuscleKg *valueobject.WeightKg
	if req.SkeletalMuscleKg != nil {
		skeletalMuscleKg, err = valueobject.NewWeightKgFromString(*req.SkeletalMuscleKg)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid skeletal muscle kg"))
			return
		}
	}
	output, err := h.record.Execute(r.Context(), weightusecase.RecordWeightInput{
		UserID:            userID,
		WeightKg:          *weightKg,
		BodyFatPercentage: bodyFatPercentage,
		SkeletalMuscleKg:  skeletalMuscleKg,
		MeasuredAt:        req.MeasuredAt,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, weightdto.RecordWeightResponse{WeightID: output.WeightID.Value()})
}
