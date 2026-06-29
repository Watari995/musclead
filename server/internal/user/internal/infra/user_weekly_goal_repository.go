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

type userWeeklyGoalRepository struct {
	dbmap *gorp.DbMap
}

func NewUserWeeklyGoalRepository(dbmap *gorp.DbMap) userdomain.UserWeeklyGoalRepository {
	return &userWeeklyGoalRepository{dbmap: dbmap}
}

const upsertUserWeeklyGoalSQL = `
INSERT INTO user_weekly_goals (id, user_id, training_count, calorie_average, weight_change_kg, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
		training_count = VALUES(training_count),
		calorie_average = VALUES(calorie_average),
		weight_change_kg = VALUES(weight_change_kg),
		updated_at = VALUES(updated_at)
`

func (r *userWeeklyGoalRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*userdomain.UserWeeklyGoal, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row UserWeeklyGoalModel
	err = q.SelectOne(&row, "SELECT id, user_id, training_count, calorie_average, weight_change_kg, created_at, updated_at FROM user_weekly_goals WHERE user_id = ?", userIDBytes)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	result, err := toUserWeeklyGoalEntity(row)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userWeeklyGoalRepository) Save(ctx context.Context, weeklyGoal *userdomain.UserWeeklyGoal) (*userdomain.UserWeeklyGoal, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertParamsFromEntity(weeklyGoal)
	if err != nil {
		return nil, err
	}
	if _, err := q.Exec(upsertUserWeeklyGoalSQL, params...); err != nil {
		return nil, err
	}
	return weeklyGoal, nil
}

func buildUpsertParamsFromEntity(weeklyGoal *userdomain.UserWeeklyGoal) ([]any, error) {
	IDBytes, err := weeklyGoal.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := weeklyGoal.UserID().Bytes()
	if err != nil {
		return nil, err
	}

	var trainingCount sql.NullInt32
	if weeklyGoal.TrainingCount() != nil {
		trainingCount = sql.NullInt32{
			Int32: int32(weeklyGoal.TrainingCount().Value()),
			Valid: true,
		}
	}
	var calorieAverage sql.NullInt32
	if weeklyGoal.CalorieAverage() != nil {
		calorieAverage = sql.NullInt32{
			Int32: int32(weeklyGoal.CalorieAverage().Value()),
			Valid: true,
		}
	}
	var weightChangeKg sql.NullString
	if weeklyGoal.WeightChangeKg() != nil {
		weightChangeKg = sql.NullString{
			String: weeklyGoal.WeightChangeKg().String(),
			Valid:  true,
		}
	}

	return []any{
		IDBytes,
		userIDBytes,
		trainingCount,
		calorieAverage,
		weightChangeKg,
		weeklyGoal.CreatedAt(),
		weeklyGoal.UpdatedAt(),
	}, nil
}

func toUserWeeklyGoalEntity(row UserWeeklyGoalModel) (*userdomain.UserWeeklyGoal, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserWeeklyGoalID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	var trainingCount *valueobject.NonNegativeInt
	if row.TrainingCount.Valid {
		trainingCount, err = valueobject.NewNonNegativeInt(int(row.TrainingCount.Int32))
		if err != nil {
			return nil, err
		}
	}
	var calorieAverage *valueobject.NonNegativeInt
	if row.CalorieAverage.Valid {
		calorieAverage, err = valueobject.NewNonNegativeInt(int(row.CalorieAverage.Int32))
		if err != nil {
			return nil, err
		}
	}
	var weightChangeKg *valueobject.WeightChangeKg
	if row.WeightChangeKg.Valid {
		weightChangeKgDecimal, err := sqlconv.DecimalFromNullString(row.WeightChangeKg)
		if err != nil {
			return nil, err
		}
		weightChangeKg, err = valueobject.NewWeightChangeKgFromDecimal(*weightChangeKgDecimal)
		if err != nil {
			return nil, err
		}
	}

	return userdomain.NewUserWeeklyGoal(
		*id,
		*userID,
		trainingCount,
		calorieAverage,
		weightChangeKg,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}
