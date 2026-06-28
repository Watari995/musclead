package calendarusecase

import (
	"context"
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetDailySummaryInput struct {
	UserID valueobject.UserID
	Date   time.Time
}

type GetDailySummaryOutput struct {
	Trainings []*trainingpublicfunctions.TrainingSummaryView
	Meals     []*mealpublicfunctions.MealSummaryView
	Weights   []*weightpublicfunctions.WeightSummaryView
}

type GetDailySummary struct {
	trainingQuery trainingpublicfunctions.TrainingQuery
	mealQuery     mealpublicfunctions.MealQuery
	weightQuery   weightpublicfunctions.WeightQuery
}

func NewGetDailySummary(
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
) *GetDailySummary {
	return &GetDailySummary{
		trainingQuery: trainingQuery,
		mealQuery:     mealQuery,
		weightQuery:   weightQuery,
	}
}

func (uc *GetDailySummary) Execute(ctx context.Context, input GetDailySummaryInput) (*GetDailySummaryOutput, error) {
	panic("not implemented")
}
