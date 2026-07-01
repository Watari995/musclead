package mealinfra

import (
	"context"
	"database/sql"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type mealQueryService struct {
	dbmap *gorp.DbMap
}

func NewMealQueryService(dbmap *gorp.DbMap) mealdomain.MealQueryService {
	return &mealQueryService{dbmap: dbmap}
}

type listMealDatesByMonthRow struct {
	MealDate time.Time `db:"meal_date"`
}

func (s *mealQueryService) ListMealDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return []time.Time{}, err
	}
	var rows []listMealDatesByMonthRow
	if _, err = q.Select(&rows, `
	SELECT DISTINCT DATE(CONVERT_TZ(eaten_at, '+00:00', '+09:00')) AS meal_date
	FROM meals
	WHERE user_id = ?
	AND YEAR(CONVERT_TZ(eaten_at, '+00:00', '+09:00')) = ?
	AND MONTH(CONVERT_TZ(eaten_at, '+00:00', '+09:00')) = ?
	ORDER BY meal_date ASC
	`, userIDBytes, year, month); err != nil {
		return []time.Time{}, err
	}
	var result []time.Time
	for _, row := range rows {
		result = append(result, row.MealDate)
	}
	return result, nil
}

type listMealSummaryByDateRow struct {
	MealID        []byte         `db:"meal_id"`
	MealType      string         `db:"meal_type"`
	EatenAt       time.Time      `db:"eaten_at"`
	Calories      int            `db:"calories"`
	ProteinG      sql.NullString `db:"protein_g"`
	FatG          sql.NullString `db:"fat_g"`
	CarbohydrateG sql.NullString `db:"carbohydrate_g"`
}

func (s *mealQueryService) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*mealdomain.MealSummaryView, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var rows []listMealSummaryByDateRow
	if _, err = q.Select(&rows, `
	SELECT id AS meal_id, meal_type, eaten_at, calories, protein_g, fat_g, carbohydrate_g
	FROM meals
	WHERE user_id = ?
	AND DATE(CONVERT_TZ(eaten_at, '+00:00', '+09:00')) = ?
	ORDER BY eaten_at ASC
	`, userIDBytes, date.Format("2006-01-02")); err != nil {
		return nil, err
	}
	var result []*mealdomain.MealSummaryView
	for _, row := range rows {
		mealSummary, err := toMealSummaryViewFromRow(row)
		if err != nil {
			return nil, err
		}
		result = append(result, mealSummary)
	}
	return result, nil
}

type getAverageCaloriesInAWeekRow struct {
	AverageCalories sql.NullFloat64 `db:"average_calories"`
}

func (s *mealQueryService) GetAverageCaloriesInAWeek(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (*valueobject.NonNegativeDecimal, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	weekStartStr := weekStart.Format("2006-01-02")
	weekEnd := weekStart.AddDate(0, 0, 6).Format("2006-01-02")
	var row getAverageCaloriesInAWeekRow
	if err := q.SelectOne(&row, `
		SELECT SUM(calories) / 7 AS average_calories
		FROM meals
		WHERE user_id = ?
		AND DATE(CONVERT_TZ(eaten_at, '+00:00', '+09:00')) BETWEEN ? AND ?
	`, userIDBytes, weekStartStr, weekEnd); err != nil {
		return nil, err
	}
	averageCalories, err := sqlconv.NewNonNegativeDecimalFromNullFloat64(row.AverageCalories)
	if err != nil {
		return nil, err
	}
	return averageCalories, nil
}

func toMealSummaryViewFromRow(row listMealSummaryByDateRow) (*mealdomain.MealSummaryView, error) {
	mealID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.MealID](row.MealID)
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

	return &mealdomain.MealSummaryView{
		MealID:        *mealID,
		MealType:      *mealType,
		EatenAt:       row.EatenAt,
		Calories:      *calories,
		ProteinG:      proteinG,
		FatG:          fatG,
		CarbohydrateG: carbohydrateG,
	}, nil
}
