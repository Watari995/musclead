package mealinfra

import (
	"context"
	"database/sql"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type mealTemplateRepository struct {
	dbmap *gorp.DbMap
}

func NewMealTemplateRepository(dbmap *gorp.DbMap) mealdomain.MealTemplateRepository {
	return &mealTemplateRepository{dbmap: dbmap}
}

const upsertMealTemplateSQL = `
INSERT INTO meal_templates (id, user_id, name, display_order, meal_type, calories, protein_g, fat_g, carbohydrate_g, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    display_order = VALUES(display_order),
    meal_type = VALUES(meal_type),
    calories = VALUES(calories),
    protein_g = VALUES(protein_g),
    fat_g = VALUES(fat_g),
    carbohydrate_g = VALUES(carbohydrate_g),
    updated_at = VALUES(updated_at)
`

func (r *mealTemplateRepository) FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*mealdomain.MealTemplate, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var mealTemplateRows []MealTemplateModel
	_, err = q.Select(&mealTemplateRows, "SELECT id, user_id, name, display_order, meal_type, calories, protein_g, fat_g, carbohydrate_g, created_at, updated_at FROM meal_templates WHERE user_id = ? ORDER BY display_order DESC, created_at DESC", bytes)
	if err != nil {
		return nil, err
	}
	result := make([]*mealdomain.MealTemplate, len(mealTemplateRows))
	for i, row := range mealTemplateRows {
		mealTemplate, err := toMealTemplate(row)
		if err != nil {
			return nil, err
		}
		result[i] = mealTemplate
	}
	return result, nil
}

func (r *mealTemplateRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*mealdomain.MealTemplate, pagination.OffsetPaginator, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	var mealTemplateRows []MealTemplateModel
	_, err = q.Select(&mealTemplateRows, "SELECT id, user_id, name, display_order, meal_type, calories, protein_g, fat_g, carbohydrate_g, created_at, updated_at FROM meal_templates WHERE user_id = ? ORDER BY display_order DESC, created_at DESC LIMIT ? OFFSET ?", bytes, int32(limit), int32(offset))
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	total, err := q.SelectInt("SELECT COUNT(*) FROM meal_templates WHERE user_id = ?", bytes)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	paginator := pagination.NewOffsetPaginator(int(total), offset, limit)
	if len(mealTemplateRows) == 0 {
		return []*mealdomain.MealTemplate{}, paginator, nil
	}
	result := make([]*mealdomain.MealTemplate, len(mealTemplateRows))
	for i, row := range mealTemplateRows {
		mealTemplate, err := toMealTemplate(row)
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		result[i] = mealTemplate
	}
	return result, paginator, nil
}

func (r *mealTemplateRepository) NextDisplayOrder(ctx context.Context, userID valueobject.UserID) (int, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return 0, err
	}
	next, err := q.SelectInt("SELECT COALESCE(MAX(display_order) + 1, 0) FROM meal_templates WHERE user_id = ?", bytes)
	if err != nil {
		return 0, err
	}
	return int(next), nil
}

func (r *mealTemplateRepository) Save(ctx context.Context, mealTemplate *mealdomain.MealTemplate) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertMealTemplateParams(mealTemplate)
	if err != nil {
		return err
	}
	_, err = q.Exec(upsertMealTemplateSQL, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r *mealTemplateRepository) DeleteByID(ctx context.Context, id valueobject.MealTemplateID) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := id.Bytes()
	if err != nil {
		return err
	}
	if _, err := q.Exec("DELETE FROM meal_templates WHERE id = ?", bytes); err != nil {
		return err
	}
	return nil
}

func toMealTemplate(row MealTemplateModel) (*mealdomain.MealTemplate, error) {
	mealTemplateID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.MealTemplateID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewString100(row.Name)
	if err != nil {
		return nil, err
	}
	displayOrder, err := valueobject.NewNonNegativeInt(int(row.DisplayOrder))
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
	return mealdomain.NewMealTemplate(*mealTemplateID, *userID, *name, *displayOrder, *mealType, *calories, proteinG, fatG, carbohydrateG, row.CreatedAt, row.UpdatedAt), nil
}

func buildUpsertMealTemplateParams(mealTemplate *mealdomain.MealTemplate) ([]any, error) {
	idBytes, err := mealTemplate.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := mealTemplate.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	name := mealTemplate.Name().Value()
	displayOrder := mealTemplate.DisplayOrder().Value()
	mealType := mealTemplate.MealType().Value()
	calories := mealTemplate.Calories().Value()
	// nullable になる可能性のあるものは値がある時だけ変換する(ないときは自動でValid falseで渡されるため変換不要なのでそのまま渡す)
	var proteinG sql.NullString
	if mealTemplate.ProteinG() != nil {
		proteinG = sqlconv.DecimalToNullString(mealTemplate.ProteinG().Value())
	}
	var fatG sql.NullString
	if mealTemplate.FatG() != nil {
		fatG = sqlconv.DecimalToNullString(mealTemplate.FatG().Value())
	}
	var carbohydrateG sql.NullString
	if mealTemplate.CarbohydrateG() != nil {
		carbohydrateG = sqlconv.DecimalToNullString(mealTemplate.CarbohydrateG().Value())
	}
	createdAt := mealTemplate.CreatedAt()
	updatedAt := mealTemplate.UpdatedAt()
	return []any{idBytes, userIDBytes, name, displayOrder, mealType, calories, proteinG, fatG, carbohydrateG, createdAt, updatedAt}, nil
}
