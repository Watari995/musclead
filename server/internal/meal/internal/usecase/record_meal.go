package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

type RecordMealInput struct {
	UserID        valueobject.UserID
	EatenAt       time.Time
	MealType      valueobject.String20
	Calories      valueobject.NonNegativeInt
	ProteinG      *valueobject.NonNegativeDecimal
	FatG          *valueobject.NonNegativeDecimal
	CarbohydrateG *valueobject.NonNegativeDecimal
	Memo          *valueobject.String1000
	FoodProductID *valueobject.FoodProductID
	ServingCount  valueobject.NonNegativeDecimal
	Photos        []mealdomain.PhotoSpec
}

type RecordMealOutput struct {
	MealID valueobject.MealID
}

type RecordMeal struct {
	mealRepo  mealdomain.MealRepository
	txManager dbtx.TransactionManager
}

func (uc *RecordMeal) Execute(ctx context.Context, input RecordMealInput) (*RecordMealOutput, error) {
	meal := mealdomain.CreateMeal(input.UserID, input.EatenAt, input.MealType, input.Calories, input.ProteinG, input.FatG, input.CarbohydrateG, input.Memo, input.FoodProductID, input.ServingCount, input.Photos)
	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		if err := uc.mealRepo.Save(txCtx, meal); err != nil {
			return myerror.NewInternalError().Wrap(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &RecordMealOutput{MealID: meal.ID()}, nil
}

func NewRecordMeal(mealRepo mealdomain.MealRepository, txManager dbtx.TransactionManager) *RecordMeal {
	return &RecordMeal{mealRepo: mealRepo, txManager: txManager}
}
