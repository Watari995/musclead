package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListMealSummaryByDateInput struct {
	UserID valueobject.UserID
	Date   time.Time
}

type ListMealSummaryByDateOutput struct {
	MealSummaries []*mealdomain.MealSummaryView
}

type ListMealSummaryByDate struct {
	mealQuery mealdomain.MealQueryService
}

func NewListMealSummaryByDate(mealQuery mealdomain.MealQueryService) *ListMealSummaryByDate {
	return &ListMealSummaryByDate{mealQuery: mealQuery}
}

func (uc *ListMealSummaryByDate) Execute(ctx context.Context, input ListMealSummaryByDateInput) (*ListMealSummaryByDateOutput, error) {
	summaries, err := uc.mealQuery.ListSummaryByDate(ctx, input.UserID, input.Date)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListMealSummaryByDateOutput{MealSummaries: summaries}, nil
}
