package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindMealByIDInput struct {
	MealID valueobject.MealID
	UserID valueobject.UserID
}

type FindMealByIDOutput struct {
	Meal *mealdomain.Meal
}

type FindMealByID struct {
	mealRepo mealdomain.MealRepository
}

func (uc *FindMealByID) Execute(ctx context.Context, input FindMealByIDInput) (*FindMealByIDOutput, error) {
	meal, err := uc.mealRepo.FindByID(ctx, input.MealID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if meal == nil {
		return nil, myerror.NewMealNotFoundError()
	}
	if meal.UserID() != input.UserID {
		return nil, myerror.NewPermissionError().SetMessage("meal does not belong to the user")
	}
	return &FindMealByIDOutput{Meal: meal}, nil
}

func NewFindMealByID(mealRepo mealdomain.MealRepository) *FindMealByID {
	return &FindMealByID{mealRepo: mealRepo}
}
