package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type RecordMealInput struct {
	UserID        valueobject.UserID
	EatenAt       time.Time
	MealType      valueobject.String20
	Calories      valueobject.NonNegativeInt
	ProteinG      valueobject.NonNegativeInt
	FatG          valueobject.NonNegativeInt
	CarbohydrateG valueobject.NonNegativeInt
	Memo          *valueobject.String1000
}

type RecordMealOutput struct {
	MealID valueobject.MealID
}

type RecordMeal struct {
	mealRepo mealdomain.MealRepository
}

func (uc *RecordMeal) Execute(ctx context.Context, input RecordMealInput) (*RecordMealOutput, error) {
	meal := mealdomain.CreateMeal(input.UserID, input.EatenAt, input.MealType, input.Calories, input.ProteinG, input.FatG, input.CarbohydrateG, input.Memo)
	if err := uc.mealRepo.Save(ctx, meal); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &RecordMealOutput{MealID: meal.ID()}, nil
}

func NewRecordMeal(mealRepo mealdomain.MealRepository) *RecordMeal {
	return &RecordMeal{mealRepo: mealRepo}
}
