package calendarhandler

import (
	"net/http"
	"strconv"
	"time"

	calendardto "github.com/Watari995/musclead/internal/calendar/dto"
	calendarusecase "github.com/Watari995/musclead/internal/calendar/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
)

type CalendarHandler struct {
	getMonthlySummary *calendarusecase.GetMonthlySummary
	getDailySummary   *calendarusecase.GetDailySummary
}

func New(
	getMonthlySummary *calendarusecase.GetMonthlySummary,
	getDailySummary *calendarusecase.GetDailySummary,
) http.Handler {
	h := &CalendarHandler{
		getMonthlySummary: getMonthlySummary,
		getDailySummary:   getDailySummary,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /calendar/monthly-summary", h.GetMonthlySummary)
	mux.HandleFunc("GET /calendar/daily-summary", h.GetDailySummary)
	return mux
}

// GetMonthlySummary godoc
//
// @Summary 月間サマリー取得
// @Tags calendar
// @Produce json
// @Security BearerAuth
// @Param year query int true "年"
// @Param month query int true "月"
// @Success 200 {object} calendardto.GetMonthlySummaryResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /calendar/monthly-summary [get]
func (h *CalendarHandler) GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid year"))
		return
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid month"))
		return
	}
	summary, err := h.getMonthlySummary.Execute(r.Context(), calendarusecase.GetMonthlySummaryInput{
		UserID: userID,
		Year:   yearInt,
		Month:  monthInt,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var res calendardto.GetMonthlySummaryResponse
	for _, daySummary := range summary.Days {
		res.Days = append(res.Days, calendardto.MonthlySummaryDayDTO{
			Date:        daySummary.Date.Format("2006-01-02"),
			HasTraining: daySummary.HasTraining,
			HasMeal:     daySummary.HasMeal,
			HasWeight:   daySummary.HasWeight,
		})
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

// GetDailySummary godoc
//
// @Summary 日別サマリー取得
// @Tags calendar
// @Produce json
// @Security BearerAuth
// @Param date query string true "日付 (YYYY-MM-DD)"
// @Success 200 {object} calendardto.GetDailySummaryResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /calendar/daily-summary [get]
func (h *CalendarHandler) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	date := r.URL.Query().Get("date")
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid date format, expected YYYY-MM-DD"))
		return
	}
	summary, err := h.getDailySummary.Execute(r.Context(), calendarusecase.GetDailySummaryInput{
		UserID: userID,
		Date:   dateTime,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	trainingDTOs := make([]calendardto.TrainingSummaryDTO, 0, len(summary.Trainings))
	for _, t := range summary.Trainings {
		trainingDTO := calendardto.FromTrainingSummaryViewToDTO(t)
		trainingDTOs = append(trainingDTOs, trainingDTO)
	}
	mealDTOs := make([]calendardto.MealSummaryDTO, 0, len(summary.Meals))
	for _, m := range summary.Meals {
		mealDTO := calendardto.FromMealSummaryViewToDTO(m)
		mealDTOs = append(mealDTOs, mealDTO)
	}
	weightDTOs := make([]calendardto.WeightSummaryDTO, 0, len(summary.Weights))
	for _, w := range summary.Weights {
		weightDTO := calendardto.FromWeightSummaryViewToDTO(w)
		weightDTOs = append(weightDTOs, weightDTO)
	}
	res := calendardto.GetDailySummaryResponse{
		Trainings:     trainingDTOs,
		Meals:         mealDTOs,
		TotalCalories: summary.TotalCalories.Value(),
		Weights:       weightDTOs,
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}
