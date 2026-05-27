package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateMealInput struct {
	MealID        valueobject.MealID
	UserID        valueobject.UserID
	EatenAt       time.Time
	MealType      valueobject.String20
	Calories      valueobject.NonNegativeInt
	ProteinG      *valueobject.NonNegativeDecimal
	FatG          *valueobject.NonNegativeDecimal
	CarbohydrateG *valueobject.NonNegativeDecimal
	Memo          *valueobject.String1000
	Photos        []mealdomain.PhotoData
}

type UpdateMealOutput struct {
	MealID valueobject.MealID
}

type UpdateMeal struct {
	mealRepo mealdomain.MealRepository
}

func (uc *UpdateMeal) Execute(ctx context.Context, input UpdateMealInput) (*UpdateMealOutput, error) {
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
	params := mealdomain.UpdateMealParams{
		EatenAt:       input.EatenAt,
		MealType:      input.MealType,
		Calories:      input.Calories,
		ProteinG:      input.ProteinG,
		FatG:          input.FatG,
		CarbohydrateG: input.CarbohydrateG,
		Memo:          input.Memo,
		Photos:        input.Photos,
	}
	meal.Update(params)
	if err := uc.mealRepo.Save(ctx, meal); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpdateMealOutput{MealID: meal.ID()}, nil
}

func NewUpdateMeal(mealRepo mealdomain.MealRepository) *UpdateMeal {
	return &UpdateMeal{mealRepo: mealRepo}
}
