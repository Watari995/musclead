package mealusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetAverageCaloriesInAWeek struct {
	mealQuery mealdomain.MealQueryService
}

func NewGetAverageCaloriesInAWeek(mealQuery mealdomain.MealQueryService) *GetAverageCaloriesInAWeek {
	return &GetAverageCaloriesInAWeek{mealQuery: mealQuery}
}

func (uc *GetAverageCaloriesInAWeek) Execute(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (*valueobject.NonNegativeDecimal, error) {
	avg, err := uc.mealQuery.GetAverageCaloriesInAWeek(ctx, userID, weekStart)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return avg, nil
}
