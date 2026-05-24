package mealinfra

import (
	"context"
	"math"

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
	bytes, err := userId.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	meals, err := r.db.FindAllMealsByUserIDWithOffsetPagination(ctx, mealdbgen.FindAllMealsByUserIDWithOffsetPaginationParams{
		UserID: bytes,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	total, err := r.db.CountMealsByUserID(ctx, bytes)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	paginator := pagination.OffsetPaginator{
		CurrentPage:  offset/limit + 1,
		ItemsPerPage: limit,
		TotalItems:   int(total),
		TotalPages:   int(math.Ceil(float64(total) / float64(limit))),
	}
	if len(meals) == 0 {
		return []*mealdomain.Meal{}, paginator, nil
	}

	photos, err := r.db.FindMealPhotosByMealIDs(ctx, lo.Map(meals, func(meal mealdbgen.Meal, _ int) []byte {
		return meal.ID
	}))
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	photosByMealID := lo.GroupBy(photos, func(photo mealdbgen.MealPhoto) string {
		return string(photo.MealID)
	})

	result := make([]*mealdomain.Meal, len(meals))
	for i, meal := range meals {
		photos := photosByMealID[string(meal.ID)]
		mealEntity, err := toMeal(meal, photos)
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		result[i] = mealEntity
	}

	return result, paginator, nil
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
