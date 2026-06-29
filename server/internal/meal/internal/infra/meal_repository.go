package mealinfra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlquery"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
	"github.com/samber/lo"
)

type mealRepository struct {
	dbmap *gorp.DbMap
}

const upsertMealSQL = `
INSERT INTO meals (id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, food_product_id, serving_count, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    eaten_at = VALUES(eaten_at),
    meal_type = VALUES(meal_type),
    calories = VALUES(calories),
    protein_g = VALUES(protein_g),
    fat_g = VALUES(fat_g),
    carbohydrate_g = VALUES(carbohydrate_g),
    memo = VALUES(memo),
    food_product_id = VALUES(food_product_id),
    serving_count = VALUES(serving_count),
    updated_at = VALUES(updated_at)
`

func (r *mealRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*mealdomain.Meal, pagination.OffsetPaginator, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	var mealRows []MealModel
	_, err = q.Select(&mealRows,
		"SELECT id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, food_product_id, serving_count, created_at, updated_at FROM meals WHERE user_id = ? ORDER BY eaten_at DESC LIMIT ? OFFSET ?",
		bytes, int32(limit), int32(offset),
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	total, err := q.SelectInt(
		"SELECT COUNT(*) FROM meals WHERE user_id = ?", bytes,
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	paginator := pagination.NewOffsetPaginator(int(total), offset, limit)

	if len(mealRows) == 0 {
		return []*mealdomain.Meal{}, paginator, nil
	}

	mealIDs := lo.Map(mealRows, func(m MealModel, _ int) []byte {
		return m.ID
	})
	photos, err := r.selectPhotosByMealIDs(q, mealIDs)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	photosByMealID := lo.GroupBy(photos, func(p MealPhotoModel) string {
		return string(p.MealID)
	})

	result := make([]*mealdomain.Meal, len(mealRows))
	for i, m := range mealRows {
		meal, err := toMeal(m, photosByMealID[string(m.ID)])
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		result[i] = meal
	}

	return result, paginator, nil
}

func (r *mealRepository) FindByIDAndUserID(ctx context.Context, id valueobject.MealID, userID valueobject.UserID) (*mealdomain.Meal, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var mealRow MealModel
	err = q.SelectOne(&mealRow,
		"SELECT id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, food_product_id, serving_count, created_at, updated_at FROM meals WHERE id = ? AND user_id = ?",
		idBytes, userIDBytes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var photos []MealPhotoModel
	_, err = q.Select(&photos,
		"SELECT id, meal_id, image_path, display_order, created_at FROM meal_photos WHERE meal_id = ? ORDER BY display_order ASC",
		idBytes,
	)
	if err != nil {
		return nil, err
	}

	return toMeal(mealRow, photos)
}

func (r *mealRepository) Save(ctx context.Context, meal *mealdomain.Meal) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := meal.ID().Bytes()
	if err != nil {
		return err
	}
	params, err := buildUpsertMealParams(meal)
	if err != nil {
		return err
	}

	if _, err := q.Exec(upsertMealSQL, params...); err != nil {
		return err
	}
	if _, err := q.Exec("DELETE FROM meal_photos WHERE meal_id = ?", bytes); err != nil {
		return err
	}
	for _, photo := range meal.Photos() {
		mealPhotoID, err := valueobject.NewPrimaryID[valueobject.MealPhotoID]().Bytes()
		if err != nil {
			return err
		}
		if _, err := q.Exec(
			"INSERT INTO meal_photos (id, meal_id, image_path, display_order, created_at) VALUES (?, ?, ?, ?, ?)",
			mealPhotoID, bytes, photo.ImagePath, photo.DisplayOrder, time.Now(),
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *mealRepository) DeleteByID(ctx context.Context, id valueobject.MealID) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := id.Bytes()
	if err != nil {
		return err
	}
	_, err = q.Exec("DELETE FROM meals WHERE id = ?", bytes)
	return err
}

// selectPhotosByMealIDs は IN 句で複数の meal_id に紐づく photos を一括取得する。
// プレースホルダ生成は shared/sqlquery を経由して同種重複を避ける。
func (r *mealRepository) selectPhotosByMealIDs(q gorp.SqlExecutor, mealIDs [][]byte) ([]MealPhotoModel, error) {
	if len(mealIDs) == 0 {
		return nil, nil
	}
	placeholders, args := sqlquery.InPlaceholders(mealIDs)
	query := fmt.Sprintf(
		"SELECT id, meal_id, image_path, display_order, created_at FROM meal_photos WHERE meal_id IN (%s) ORDER BY meal_id ASC, display_order ASC",
		placeholders,
	)
	var photos []MealPhotoModel
	_, err := q.Select(&photos, query, args...)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func toPhotoSpec(photos []MealPhotoModel) []mealdomain.PhotoSpec {
	return lo.Map(photos, func(p MealPhotoModel, _ int) mealdomain.PhotoSpec {
		return mealdomain.PhotoSpec{
			ImagePath:    p.ImagePath,
			DisplayOrder: p.DisplayOrder,
		}
	})
}

func toMeal(row MealModel, photos []MealPhotoModel) (*mealdomain.Meal, error) {
	mealID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.MealID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	mealType, err := valueobject.NewString20(row.MealType)
	if err != nil {
		return nil, err
	}
	calories, err := valueobject.NewNonNegativeInt(row.Calories)
	if err != nil {
		return nil, err
	}
	proteinG, err := sqlconv.NewNonNegativeDecimalFromNullString(row.ProteinG)
	if err != nil {
		return nil, err
	}
	fatG, err := sqlconv.NewNonNegativeDecimalFromNullString(row.FatG)
	if err != nil {
		return nil, err
	}
	carbohydrateG, err := sqlconv.NewNonNegativeDecimalFromNullString(row.CarbohydrateG)
	if err != nil {
		return nil, err
	}
	var memoVO *valueobject.String1000
	if row.Memo.Valid {
		memoVO, err = valueobject.NewString1000(row.Memo.String)
		if err != nil {
			return nil, err
		}
	}
	foodProductID, err := sqlconv.NewPrimaryIDFromNullableBytes[valueobject.FoodProductID](row.FoodProductID)
	if err != nil {
		return nil, err
	}
	servingCount, err := valueobject.NewNonNegativeDecimalFromString(row.ServingCount)
	if err != nil {
		return nil, err
	}
	return mealdomain.NewMeal(*mealID, *userID, row.EatenAt, *mealType, *calories, proteinG, fatG, carbohydrateG, memoVO, foodProductID, *servingCount, row.CreatedAt, row.UpdatedAt, toPhotoSpec(photos)), nil
}

func buildUpsertMealParams(meal *mealdomain.Meal) ([]any, error) {
	bytes, err := meal.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := meal.UserID().Bytes()
	if err != nil {
		return nil, err
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
	foodProductIDBytes, err := sqlconv.NewBytesFromNullablePrimaryID(meal.FoodProductID())
	if err != nil {
		return nil, err
	}
	return []any{
		bytes,
		userIDBytes,
		meal.EatenAt(),
		meal.MealType().Value(),
		meal.Calories().Value(),
		proteinG,
		fatG,
		carbohydrateG,
		memo,
		foodProductIDBytes,
		meal.ServingCount().Value().String(),
		meal.CreatedAt(),
		meal.UpdatedAt(),
	}, nil
}

func NewMealRepository(dbmap *gorp.DbMap) mealdomain.MealRepository {
	return &mealRepository{dbmap: dbmap}
}
