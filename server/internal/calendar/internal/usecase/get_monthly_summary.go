package calendarusecase

import (
	"context"
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetMonthlySummaryInput struct {
	UserID valueobject.UserID
	Year   int
	Month  int
}

type DaySummary struct {
	Date        time.Time
	HasTraining bool
	HasMeal     bool
	HasWeight   bool
}

type GetMonthlySummaryOutput struct {
	Days []DaySummary
}

type GetMonthlySummary struct {
	trainingQuery trainingpublicfunctions.TrainingQuery
	mealQuery     mealpublicfunctions.MealQuery
	weightQuery   weightpublicfunctions.WeightQuery
}

func NewGetMonthlySummary(
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
) *GetMonthlySummary {
	return &GetMonthlySummary{
		trainingQuery: trainingQuery,
		mealQuery:     mealQuery,
		weightQuery:   weightQuery,
	}
}

func (uc *GetMonthlySummary) Execute(ctx context.Context, input GetMonthlySummaryInput) (*GetMonthlySummaryOutput, error) {
	panic("not implemented")
}
