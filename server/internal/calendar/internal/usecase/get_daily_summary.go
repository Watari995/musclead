package calendarusecase

import (
	"context"
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	"golang.org/x/sync/errgroup"
)

type GetDailySummaryInput struct {
	UserID valueobject.UserID
	Date   time.Time
}

type GetDailySummaryOutput struct {
	Trainings     []*trainingpublicfunctions.TrainingSummaryView
	Meals         []*mealpublicfunctions.MealSummaryView
	TotalCalories valueobject.NonNegativeInt
	Weights       []*weightpublicfunctions.WeightSummaryView
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
	g, ctx := errgroup.WithContext(ctx)

	var trainings []*trainingpublicfunctions.TrainingSummaryView
	g.Go(func() error {
		var err error
		trainings, err = uc.trainingQuery.ListSummaryByDate(ctx, input.UserID, input.Date)
		return err
	})

	var meals []*mealpublicfunctions.MealSummaryView
	g.Go(func() error {
		var err error
		meals, err = uc.mealQuery.ListSummaryByDate(ctx, input.UserID, input.Date)
		return err
	})

	var weights []*weightpublicfunctions.WeightSummaryView
	g.Go(func() error {
		var err error
		weights, err = uc.weightQuery.ListSummaryByDate(ctx, input.UserID, input.Date)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	totalCalories := valueobject.NonNegativeInt{}
	for _, meal := range meals {
		totalCalories = totalCalories.Add(meal.Calories)
	}

	return &GetDailySummaryOutput{
		Trainings:     trainings,
		Meals:         meals,
		TotalCalories: totalCalories,
		Weights:       weights,
	}, nil
}
