package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateMealInput struct {
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

type CreateMealOutput struct {
	MealID valueobject.MealID
}

type CreateMeal struct {
	mealRepo mealdomain.MealRepository
}

func (uc *CreateMeal) Execute(ctx context.Context, input CreateMealInput) (*CreateMealOutput, error) {
	meal := mealdomain.CreateMeal(
		input.UserID,
		input.EatenAt,
		input.MealType,
		input.Calories,
		input.ProteinG,
		input.FatG,
		input.CarbohydrateG,
		input.Memo,
		input.Photos,
	)
	if err := uc.mealRepo.Save(ctx, meal); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &CreateMealOutput{MealID: meal.ID()}, nil
}

func NewCreateMeal(mealRepo mealdomain.MealRepository) *CreateMeal {
	return &CreateMeal{mealRepo: mealRepo}
}
