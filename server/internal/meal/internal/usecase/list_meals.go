package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListMealsInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListMealsOutput struct {
	Meals      []*mealdomain.Meal
	Pagination pagination.OffsetPaginator
}

type ListMeals struct {
	mealRepo mealdomain.MealRepository
}

func (uc *ListMeals) Execute(ctx context.Context, input ListMealsInput) (*ListMealsOutput, error) {
	meals, paginator, err := uc.mealRepo.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListMealsOutput{Meals: meals, Pagination: paginator}, nil
}

func NewListMeals(mealRepo mealdomain.MealRepository) *ListMeals {
	return &ListMeals{mealRepo: mealRepo}
}
