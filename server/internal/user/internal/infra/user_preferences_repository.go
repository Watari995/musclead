package userinfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type userPreferencesRepository struct {
	dbmap *gorp.DbMap
}

const upsertUserPreferencesSQL = `
INSERT INTO user_preferences (id, user_id, theme, meal_color, training_color, weight_color, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    theme = VALUES(theme),
    meal_color = VALUES(meal_color),
    training_color = VALUES(training_color),
    weight_color = VALUES(weight_color),
    updated_at = VALUES(updated_at)
`

func (r *userPreferencesRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*userdomain.UserPreferences, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row UserPreferencesModel
	err = q.SelectOne(&row, "SELECT id, user_id, theme, meal_color, training_color, weight_color, created_at, updated_at FROM user_preferences WHERE user_id = ?", bytes)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toUserPreferences(row)
}

func (r *userPreferencesRepository) Save(ctx context.Context, preferences *userdomain.UserPreferences) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := preferences.ID().Bytes()
	if err != nil {
		return err
	}
	userIDBytes, err := preferences.UserID().Bytes()
	if err != nil {
		return err
	}
	theme := preferences.Theme().Value()
	mealColor := preferences.MealColor().Value()
	trainingColor := preferences.TrainingColor().Value()
	weightColor := preferences.WeightColor().Value()
	_, err = q.Exec(upsertUserPreferencesSQL, bytes, userIDBytes, theme, mealColor, trainingColor, weightColor, preferences.CreatedAt(), preferences.UpdatedAt())
	return err
}

func toUserPreferences(row UserPreferencesModel) (*userdomain.UserPreferences, error) {
	userPreferencesID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserPreferencesID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	theme, err := valueobject.NewTheme(row.Theme)
	if err != nil {
		return nil, err
	}
	mealColor, err := valueobject.NewColorHex(row.MealColor)
	if err != nil {
		return nil, err
	}
	trainingColor, err := valueobject.NewColorHex(row.TrainingColor)
	if err != nil {
		return nil, err
	}
	weightColor, err := valueobject.NewColorHex(row.WeightColor)
	if err != nil {
		return nil, err
	}
	return userdomain.NewUserPreferences(
		*userPreferencesID,
		*userID,
		*theme,
		*mealColor,
		*trainingColor,
		*weightColor,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}

func NewUserPreferencesRepository(dbmap *gorp.DbMap) userdomain.UserPreferencesRepository {
	return &userPreferencesRepository{dbmap: dbmap}
}
