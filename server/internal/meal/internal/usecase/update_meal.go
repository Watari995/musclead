package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
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
	Photos        []mealdomain.PhotoSpec
}

type UpdateMealOutput struct {
	MealID valueobject.MealID
}

type UpdateMeal struct {
	mealRepo  mealdomain.MealRepository
	txManager dbtx.TransactionManager
}

func (uc *UpdateMeal) Execute(ctx context.Context, input UpdateMealInput) (*UpdateMealOutput, error) {
	var mealID valueobject.MealID
	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		meal, err := uc.mealRepo.FindByID(txCtx, input.MealID)
		if err != nil {
			return myerror.NewInternalError().Wrap(err)
		}
		if meal == nil {
			return myerror.NewMealNotFoundError()
		}
		if meal.UserID() != input.UserID {
			return myerror.NewPermissionError().SetMessage("meal does not belong to the user")
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
		if err := uc.mealRepo.Save(txCtx, meal); err != nil {
			return myerror.NewInternalError().Wrap(err)
		}
		mealID = meal.ID()
		return nil
	}); err != nil {
		return nil, err
	}
	return &UpdateMealOutput{MealID: mealID}, nil
}

func NewUpdateMeal(mealRepo mealdomain.MealRepository, txManager dbtx.TransactionManager) *UpdateMeal {
	return &UpdateMeal{mealRepo: mealRepo, txManager: txManager}
}
