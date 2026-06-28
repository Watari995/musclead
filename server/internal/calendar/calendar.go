// Package calendar は calendar モジュールの公開 Facade。
// training / meal / weight モジュールの公開インターフェース経由でデータを集約する。
package calendar

import (
	"net/http"

	calendarhandler "github.com/Watari995/musclead/internal/calendar/internal/handler"
	calendarusecase "github.com/Watari995/musclead/internal/calendar/internal/usecase"
	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

type Module struct {
	Handler http.Handler
}

func NewModule(
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
) *Module {
	getMonthlySummary := calendarusecase.NewGetMonthlySummary(trainingQuery, mealQuery, weightQuery)
	getDailySummary := calendarusecase.NewGetDailySummary(trainingQuery, mealQuery, weightQuery)

	return &Module{
		Handler: calendarhandler.New(getMonthlySummary, getDailySummary),
	}
}
