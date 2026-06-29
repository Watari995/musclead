package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListMealDatesByMonthInput struct {
	UserID valueobject.UserID
	Year   int
	Month  int
}

type ListMealDatesByMonthOutput struct {
	Dates []time.Time
}

type ListMealDatesByMonth struct {
	mealQuery mealdomain.MealQueryService
}

func NewListMealDatesByMonth(mealQuery mealdomain.MealQueryService) *ListMealDatesByMonth {
	return &ListMealDatesByMonth{mealQuery: mealQuery}
}

func (uc *ListMealDatesByMonth) Execute(ctx context.Context, input ListMealDatesByMonthInput) (*ListMealDatesByMonthOutput, error) {
	dates, err := uc.mealQuery.ListMealDatesByMonth(ctx, input.UserID, input.Year, input.Month)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListMealDatesByMonthOutput{Dates: dates}, nil
}
