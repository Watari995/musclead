package mealinfra

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	mealdbgen "github.com/Watari995/musclead/internal/meal/internal/infra/dbgen"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
)

type mealRepository struct {
	db      *sql.DB
	queries *mealdbgen.Queries
}

func (r *mealRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userId valueobject.UserID, limit int, offset int) ([]*mealdomain.Meal, pagination.OffsetPaginator, error) {
	bytes, err := userId.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	meals, err := r.queries.FindAllMealsByUserIDWithOffsetPagination(ctx, mealdbgen.FindAllMealsByUserIDWithOffsetPaginationParams{
		UserID: bytes,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	total, err := r.queries.CountMealsByUserID(ctx, bytes)
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

	photos, err := r.queries.FindMealPhotosByMealIDs(ctx, lo.Map(meals, func(meal mealdbgen.Meal, _ int) []byte {
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

func (r *mealRepository) FindByID(ctx context.Context, id valueobject.MealID) (*mealdomain.Meal, error) {
	bytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	meal, err := r.queries.FindMealByID(ctx, bytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	photos, err := r.queries.FindMealPhotosByMealID(ctx, bytes)
	if err != nil {
		return nil, err
	}
	return toMeal(meal, photos)
}

func (r *mealRepository) Save(ctx context.Context, meal *mealdomain.Meal) error {
	bytes, err := meal.ID().Bytes()
	if err != nil {
		return err
	}
	params, err := toUpsertMealParams(meal)
	if err != nil {
		return err
	}
	return dbtx.RunInTx(ctx, r.db, func(tx *sql.Tx) error {
		q := r.queries.WithTx(tx)
		if err := q.UpsertMeal(ctx, params); err != nil {
			return err
		}
		if err := q.DeleteMealPhotosByMealID(ctx, bytes); err != nil {
			return err
		}
		for _, photo := range meal.Photos() {
			mealPhotoId := valueobject.NewPrimaryId[valueobject.MealPhotoID]()
			mealPhotoIdBytes, err := mealPhotoId.Bytes()
			if err != nil {
				return err
			}
			if err := q.CreateMealPhoto(ctx, mealdbgen.CreateMealPhotoParams{
				ID:           mealPhotoIdBytes,
				MealID:       bytes,
				ImagePath:    photo.ImagePath,
				DisplayOrder: int32(photo.DisplayOrder),
				CreatedAt:    time.Now(),
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *mealRepository) DeleteByID(ctx context.Context, id valueobject.MealID) error {
	bytes, err := id.Bytes()
	if err != nil {
		return err
	}
	return r.queries.DeleteMealByID(ctx, bytes)
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

func toUpsertMealParams(meal *mealdomain.Meal) (mealdbgen.UpsertMealParams, error) {
	bytes, err := meal.ID().Bytes()
	if err != nil {
		return mealdbgen.UpsertMealParams{}, err
	}
	userIdBytes, err := meal.UserID().Bytes()
	if err != nil {
		return mealdbgen.UpsertMealParams{}, err
	}
	var proteinG sql.NullString
	if meal.ProteinG() != nil {
		proteinG = sqlconv.DecimalToNullString(meal.ProteinG().Value())
	}
	var fatG sql.NullString
	if meal.FatG() != nil {
		fatG = sqlconv.DecimalToNullString(meal.FatG().Value())
	}
	var carbohydrateG sql.NullString
	if meal.CarbohydrateG() != nil {
		carbohydrateG = sqlconv.DecimalToNullString(meal.CarbohydrateG().Value())
	}
	var memo sql.NullString
	if meal.Memo() != nil {
		memo = sqlconv.StringToNullString(meal.Memo().Value())
	}
	return mealdbgen.UpsertMealParams{
		ID:            bytes,
		UserID:        userIdBytes,
		EatenAt:       meal.EatenAt(),
		MealType:      meal.MealType().Value(),
		Calories:      int32(meal.Calories().Value()),
		ProteinG:      proteinG,
		FatG:          fatG,
		CarbohydrateG: carbohydrateG,
		Memo:          memo,
		CreatedAt:     meal.CreatedAt(),
		UpdatedAt:     meal.UpdatedAt(),
	}, nil
}

func NewMealRepository(db *sql.DB) mealdomain.MealRepository {
	return &mealRepository{db: db, queries: mealdbgen.New(db)}
}
