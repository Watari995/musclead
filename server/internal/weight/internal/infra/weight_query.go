package weightinfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
	"github.com/go-gorp/gorp/v3"
)

type weightQueryService struct {
	dbmap *gorp.DbMap
}

func NewWeightQueryService(dbmap *gorp.DbMap) weightdomain.WeightQueryService {
	return &weightQueryService{dbmap: dbmap}
}

type listWeightDatesByMonthRow struct {
	WeightDate time.Time `db:"weight_date"`
}

func (s *weightQueryService) ListWeightDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return []time.Time{}, err
	}
	var rows []listWeightDatesByMonthRow
	if _, err = q.Select(&rows, `
	SELECT DISTINCT DATE(CONVERT_TZ(measured_at, '+00:00', '+09:00')) AS weight_date
	FROM weights
	WHERE user_id = ?
	AND YEAR(CONVERT_TZ(measured_at, '+00:00', '+09:00')) = ?
	AND MONTH(CONVERT_TZ(measured_at, '+00:00', '+09:00')) = ?
	ORDER BY weight_date ASC
	`, userIDBytes, year, month); err != nil {
		return []time.Time{}, err
	}
	var result []time.Time
	for _, row := range rows {
		result = append(result, row.WeightDate)
	}
	return result, nil
}

type listWeightSummaryByDateRow struct {
	WeightID          []byte         `db:"weight_id"`
	WeightKg          string         `db:"weight_kg"`
	BodyFatPercentage sql.NullString `db:"body_fat_percentage"`
	SkeletalMuscleKg  sql.NullString `db:"skeletal_muscle_kg"`
	MeasuredAt        time.Time      `db:"measured_at"`
}

func (s *weightQueryService) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*weightdomain.WeightSummaryView, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}

	var rows []listWeightSummaryByDateRow
	if _, err = q.Select(&rows, `
	SELECT id AS weight_id, weight_kg, body_fat_percentage, skeletal_muscle_kg, measured_at
	FROM weights
	WHERE user_id = ?
	AND DATE(CONVERT_TZ(measured_at, '+00:00', '+09:00')) = ?
	ORDER BY measured_at DESC
	`, userIDBytes, date.Format("2006-01-02")); err != nil {
		return nil, err
	}

	result := make([]*weightdomain.WeightSummaryView, 0, len(rows))
	for _, row := range rows {
		weightSummary, err := toWeightSummaryViewFromRow(row)
		if err != nil {
			return nil, err
		}
		result = append(result, weightSummary)
	}
	return result, nil
}

func toWeightSummaryViewFromRow(row listWeightSummaryByDateRow) (*weightdomain.WeightSummaryView, error) {
	weightID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.WeightID](row.WeightID)
	if err != nil {
		return nil, err
	}
	weightKg, err := valueobject.NewWeightKgFromString(row.WeightKg)
	if err != nil {
		return nil, err
	}
	var bodyFatPercentage *valueobject.Percentage
	if row.BodyFatPercentage.Valid {
		bfp, err := valueobject.NewPercentageFromString(row.BodyFatPercentage.String)
		if err != nil {
			return nil, err
		}
		bodyFatPercentage = bfp
	}
	var skeletalMuscleKg *valueobject.WeightKg
	if row.SkeletalMuscleKg.Valid {
		smk, err := valueobject.NewWeightKgFromString(row.SkeletalMuscleKg.String)
		if err != nil {
			return nil, err
		}
		skeletalMuscleKg = smk
	}

	return &weightdomain.WeightSummaryView{
		WeightID:          *weightID,
		WeightKg:          *weightKg,
		BodyFatPercentage: bodyFatPercentage,
		SkeletalMuscleKg:  skeletalMuscleKg,
		MeasuredAt:        row.MeasuredAt,
	}, nil
}
