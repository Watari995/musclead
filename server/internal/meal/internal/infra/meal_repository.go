package mealinfra

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	mealdbgen "github.com/Watari995/musclead/internal/meal/internal/infra/dbgen"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
)

type mealRepository struct {
	db *mealdbgen.Queries
}

func (r *mealRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userId valueobject.UserID, limit int, offset int) ([]*mealdomain.Meal, pagination.OffsetPaginator, error) {
	// TODO: implement
	return nil, pagination.OffsetPaginator{}, nil
}

func toPhotoData(photos []mealdbgen.MealPhoto) []mealdomain.PhotoData {
	return lo.Map(photos, func(photo mealdbgen.MealPhoto, _ int) mealdomain.PhotoData {
		return mealdomain.PhotoData{
			ImagePath:    photo.ImagePath,
			DisplayOrder: int(photo.DisplayOrder),
		}
	})
}

func toMeal(mealRow mealdbgen.Meal, photos []mealdbgen.MealPhoto) (*mealdomain.Meal, error) {
	mealIdString, err := sqlconv.UUIDStringFromBytes(mealRow.ID)
	if err != nil {
		return nil, err
	}
	mealId, err := valueobject.NewPrimaryIdFromString[valueobject.MealID](mealIdString)
	if err != nil {
		return nil, err
	}
	userIdString, err := sqlconv.UUIDStringFromBytes(mealRow.UserID)
	if err != nil {
		return nil, err
	}
	userId, err := valueobject.NewPrimaryIdFromString[valueobject.UserID](userIdString)
	if err != nil {
		return nil, err
	}
	eatenAt := mealRow.EatenAt
	mealType, err := valueobject.NewString20(mealRow.MealType)
	if err != nil {
		return nil, err
	}
	calories, err := valueobject.NewNonNegativeInt(int(mealRow.Calories))
	if err != nil {
		return nil, err
	}
	proteinG, err := sqlconv.NewNonNegativeDecimalFromNullString(mealRow.ProteinG)
	if err != nil {
		return nil, err
	}
	fatG, err := sqlconv.NewNonNegativeDecimalFromNullString(mealRow.FatG)
	if err != nil {
		return nil, err
	}
	carbohydrateG, err := sqlconv.NewNonNegativeDecimalFromNullString(mealRow.CarbohydrateG)
	if err != nil {
		return nil, err
	}
	var memoVO *valueobject.String1000
	if mealRow.Memo.Valid {
		memoVO, err = valueobject.NewString1000(mealRow.Memo.String)
		if err != nil {
			return nil, err
		}
	}
	return mealdomain.NewMeal(*mealId, *userId, eatenAt, *mealType, *calories, proteinG, fatG, carbohydrateG, memoVO, mealRow.CreatedAt, mealRow.UpdatedAt, toPhotoData(photos)), nil
}
