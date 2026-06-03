package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type DeleteMealByIDInput struct {
	MealID valueobject.MealID
	UserID valueobject.UserID
}

type DeleteMealByID struct {
	mealRepo mealdomain.MealRepository
}

func (uc *DeleteMealByID) Execute(ctx context.Context, input DeleteMealByIDInput) error {
	meal, err := uc.mealRepo.FindByIDAndUserID(ctx, input.MealID, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if meal == nil {
		return myerror.NewMealNotFoundError().SetMessage("meal not found")
	}
	if err := uc.mealRepo.DeleteByID(ctx, input.MealID); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewDeleteMealByID(mealRepo mealdomain.MealRepository) *DeleteMealByID {
	return &DeleteMealByID{mealRepo: mealRepo}
}
