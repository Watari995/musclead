package mealusecase

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
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
	mealRepo      mealdomain.MealRepository
	storageClient shareddomain.StorageClient
	txManager     dbtx.TransactionManager
}

func (uc *UpdateMeal) Execute(ctx context.Context, input UpdateMealInput) (*UpdateMealOutput, error) {
	var mealID valueobject.MealID
	oldPaths := make([]string, 0)
	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		meal, err := uc.mealRepo.FindByIDAndUserID(txCtx, input.MealID, input.UserID)
		if err != nil {
			return myerror.NewInternalError().Wrap(err)
		}
		if meal == nil {
			return myerror.NewMealNotFoundError()
		}

		// old pathを全て取得する
		for _, p := range meal.Photos() {
			oldPaths = append(oldPaths, p.ImagePath)
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
	// txの後でdeleteする
	cleanUpRemovedPhotos(ctx, uc.storageClient, oldPaths, input.Photos)
	return &UpdateMealOutput{MealID: mealID}, nil
}

// 差分を取得してdeleteする
func cleanUpRemovedPhotos(
	ctx context.Context,
	storageClient shareddomain.StorageClient,
	oldPaths []string,
	newPhotos []mealdomain.PhotoSpec,
) {
	newPathSet := make(map[string]bool, len(newPhotos))
	for _, p := range newPhotos {
		newPathSet[p.ImagePath] = true
	}
	for _, o := range oldPaths {
		if !newPathSet[o] {
			// best effort
			_ = storageClient.DeleteObject(ctx, o)
		}
	}
}

func NewUpdateMeal(mealRepo mealdomain.MealRepository, txManager dbtx.TransactionManager, storageClient shareddomain.StorageClient) *UpdateMeal {
	return &UpdateMeal{mealRepo: mealRepo, txManager: txManager, storageClient: storageClient}
}
