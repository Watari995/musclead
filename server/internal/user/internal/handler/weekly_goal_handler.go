package userhandler

import (
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	userdto "github.com/Watari995/musclead/internal/user/dto"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/shopspring/decimal"
)

type WeeklyGoalHandler struct {
	getWeeklyGoal    *userusecase.GetWeeklyGoal
	upsertWeeklyGoal *userusecase.UpsertWeeklyGoal
}

func RegisterAuthenticatedWeeklyGoalHandlers(mux *http.ServeMux, getWeeklyGoal *userusecase.GetWeeklyGoal, upsertWeeklyGoal *userusecase.UpsertWeeklyGoal) {
	h := &WeeklyGoalHandler{getWeeklyGoal: getWeeklyGoal, upsertWeeklyGoal: upsertWeeklyGoal}
	mux.HandleFunc("GET /users/me/weekly-goal", h.GetWeeklyGoal)
	mux.HandleFunc("PUT /users/me/weekly-goal", h.UpsertWeeklyGoal)
}

// GetWeeklyGoal godoc
//
// @Summary 週次目標取得
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userdto.WeeklyGoalDTO
// @Failure 401 {object} httpx.ErrorResponse
// @Router /users/me/weekly-goal [get]
func (h *WeeklyGoalHandler) GetWeeklyGoal(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	output, err := h.getWeeklyGoal.Execute(r.Context(), userusecase.GetWeeklyGoalInput{UserID: userID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, userdto.WeeklyGoalFromEntity(output.Goal))
}

// UpsertWeeklyGoal godoc
//
// @Summary 週次目標更新
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userdto.UpsertWeeklyGoalRequest true "request"
// @Success 200 {object} userdto.WeeklyGoalDTO
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /users/me/weekly-goal [put]
func (h *WeeklyGoalHandler) UpsertWeeklyGoal(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req userdto.UpsertWeeklyGoalRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}

	var trainingCount *valueobject.NonNegativeInt
	if req.TrainingCount != nil {
		trainingCount, err = valueobject.NewNonNegativeInt(*req.TrainingCount)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid training_count"))
			return
		}
	}
	var calorieAverage *valueobject.NonNegativeInt
	if req.CalorieAverage != nil {
		calorieAverage, err = valueobject.NewNonNegativeInt(*req.CalorieAverage)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid calorie_average"))
			return
		}
	}
	var weightChangeKg *valueobject.WeightChangeKg
	if req.WeightChangeKg != nil {
		weightChangeKg, err = valueobject.NewWeightChangeKgFromDecimal(decimal.NewFromFloat(*req.WeightChangeKg))
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid weight_change_kg"))
			return
		}
	}

	output, err := h.upsertWeeklyGoal.Execute(r.Context(), userusecase.UpsertWeeklyGoalInput{
		UserID:         userID,
		TrainingCount:  trainingCount,
		CalorieAverage: calorieAverage,
		WeightChangeKg: weightChangeKg,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, userdto.WeeklyGoalFromEntity(output.Goal))
}
