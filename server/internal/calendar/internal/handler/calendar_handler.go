package calendarhandler

import (
	"net/http"

	calendarusecase "github.com/Watari995/musclead/internal/calendar/internal/usecase"
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

func (h *CalendarHandler) GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *CalendarHandler) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
