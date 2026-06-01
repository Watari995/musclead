package mealinfra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/pagination"
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
INSERT INTO meals (id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    eaten_at = VALUES(eaten_at),
    meal_type = VALUES(meal_type),
    calories = VALUES(calories),
    protein_g = VALUES(protein_g),
    fat_g = VALUES(fat_g),
    carbohydrate_g = VALUES(carbohydrate_g),
    memo = VALUES(memo),
    updated_at = VALUES(updated_at)
`

func (r *mealRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*mealdomain.Meal, pagination.OffsetPaginator, error) {
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	var mealRows []MealModel
	_, err = r.dbmap.WithContext(ctx).Select(&mealRows,
		"SELECT id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, created_at, updated_at FROM meals WHERE user_id = ? ORDER BY eaten_at DESC LIMIT ? OFFSET ?",
		bytes, int32(limit), int32(offset),
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	total, err := r.dbmap.WithContext(ctx).SelectInt(
		"SELECT COUNT(*) FROM meals WHERE user_id = ?", bytes,
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	paginator := pagination.OffsetPaginator{
		CurrentPage:  offset/limit + 1,
		ItemsPerPage: limit,
		TotalItems:   int(total),
		TotalPages:   int(math.Ceil(float64(total) / float64(limit))),
	}
	if len(mealRows) == 0 {
		return []*mealdomain.Meal{}, paginator, nil
	}

	mealIDs := lo.Map(mealRows, func(m MealModel, _ int) []byte {
		return m.ID
	})
	photos, err := r.selectPhotosByMealIDs(ctx, r.dbmap, mealIDs)
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

func (r *mealRepository) FindByID(ctx context.Context, id valueobject.MealID) (*mealdomain.Meal, error) {
	bytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	var mealRow MealModel
	err = r.dbmap.WithContext(ctx).SelectOne(&mealRow,
		"SELECT id, user_id, eaten_at, meal_type, calories, protein_g, fat_g, carbohydrate_g, memo, created_at, updated_at FROM meals WHERE id = ?",
		bytes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var photos []MealPhotoModel
	_, err = r.dbmap.WithContext(ctx).Select(&photos,
		"SELECT id, meal_id, image_path, display_order, created_at FROM meal_photos WHERE meal_id = ? ORDER BY display_order ASC",
		bytes,
	)
	if err != nil {
		return nil, err
	}

	return toMeal(mealRow, photos)
}

func (r *mealRepository) Save(ctx context.Context, meal *mealdomain.Meal) error {
	bytes, err := meal.ID().Bytes()
	if err != nil {
		return err
	}
	params, err := buildUpsertMealParams(meal)
	if err != nil {
		return err
	}

	tx, err := r.dbmap.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txCtx := tx.WithContext(ctx)

	if _, err := txCtx.Exec(upsertMealSQL, params...); err != nil {
		return err
	}
	if _, err := txCtx.Exec("DELETE FROM meal_photos WHERE meal_id = ?", bytes); err != nil {
		return err
	}
	for _, photo := range meal.Photos() {
		mealPhotoID, err := valueobject.NewPrimaryID[valueobject.MealPhotoID]().Bytes()
		if err != nil {
			return err
		}
		if _, err := txCtx.Exec(
			"INSERT INTO meal_photos (id, meal_id, image_path, display_order, created_at) VALUES (?, ?, ?, ?, ?)",
			mealPhotoID, bytes, photo.ImagePath, int32(photo.DisplayOrder), time.Now(),
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *mealRepository) DeleteByID(ctx context.Context, id valueobject.MealID) error {
	bytes, err := id.Bytes()
	if err != nil {
		return err
	}
	_, err = r.dbmap.WithContext(ctx).Exec("DELETE FROM meals WHERE id = ?", bytes)
	return err
}

// selectPhotosByMealIDs は IN 句で複数の meal_id に紐づく photos を一括取得する。
// プレースホルダ生成は shared/sqlquery を経由して同種重複を避ける。
func (r *mealRepository) selectPhotosByMealIDs(ctx context.Context, dbmap *gorp.DbMap, mealIDs [][]byte) ([]MealPhotoModel, error) {
	if len(mealIDs) == 0 {
		return nil, nil
	}
	placeholders, args := sqlquery.InPlaceholders(mealIDs)
	query := fmt.Sprintf(
		"SELECT id, meal_id, image_path, display_order, created_at FROM meal_photos WHERE meal_id IN (%s) ORDER BY meal_id ASC, display_order ASC",
		placeholders,
	)
	var photos []MealPhotoModel
	_, err := dbmap.WithContext(ctx).Select(&photos, query, args...)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func toPhotoData(photos []MealPhotoModel) []mealdomain.PhotoData {
	return lo.Map(photos, func(p MealPhotoModel, _ int) mealdomain.PhotoData {
		return mealdomain.PhotoData{
			ImagePath:    p.ImagePath,
			DisplayOrder: int(p.DisplayOrder),
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
	calories, err := valueobject.NewNonNegativeInt(int(row.Calories))
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
	return mealdomain.NewMeal(*mealID, *userID, row.EatenAt, *mealType, *calories, proteinG, fatG, carbohydrateG, memoVO, row.CreatedAt, row.UpdatedAt, toPhotoData(photos)), nil
}

func buildUpsertMealParams(meal *mealdomain.Meal) ([]interface{}, error) {
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
	return []interface{}{
		bytes,
		userIDBytes,
		meal.EatenAt(),
		meal.MealType().Value(),
		int32(meal.Calories().Value()),
		proteinG,
		fatG,
		carbohydrateG,
		memo,
		meal.CreatedAt(),
		meal.UpdatedAt(),
	}, nil
}

func NewMealRepository(dbmap *gorp.DbMap) mealdomain.MealRepository {
	return &mealRepository{dbmap: dbmap}
}
