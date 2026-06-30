package traininginfra

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type trainingQueryService struct {
	dbmap *gorp.DbMap
}

func NewTrainingQueryService(dbmap *gorp.DbMap) trainingdomain.TrainingQueryService {
	return &trainingQueryService{dbmap: dbmap}
}

type listTrainingDatesByMonthRow struct {
	TrainingDate time.Time `db:"training_date"`
}

func (s *trainingQueryService) ListTrainingDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return []time.Time{}, err
	}
	var rows []listTrainingDatesByMonthRow
	if _, err = q.Select(&rows, `
	SELECT DISTINCT DATE(CONVERT_TZ(started_at, '+00:00', '+09:00')) AS training_date
	FROM trainings
	WHERE user_id = ?
	AND YEAR(CONVERT_TZ(started_at, '+00:00', '+09:00')) = ?
	AND MONTH(CONVERT_TZ(started_at, '+00:00', '+09:00')) = ?
	AND EXISTS (
		SELECT 1
		FROM training_exercises te
		JOIN training_sets ts ON ts.training_exercise_id = te.id
		WHERE te.training_id = trainings.id
	)
	ORDER BY training_date ASC
	`, userIDBytes, year, month); err != nil {
		return []time.Time{}, err
	}
	var result []time.Time
	for _, row := range rows {
		result = append(result, row.TrainingDate)
	}
	return result, nil
}

type listTrainingSummaryByDateRow struct {
	TrainingID    []byte     `db:"training_id"`
	StartedAt     time.Time  `db:"started_at"`
	EndedAt       *time.Time `db:"ended_at"`
	ExerciseCount int        `db:"exercise_count"`
	SetCount      int        `db:"set_count"`
}

func (s *trainingQueryService) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*trainingdomain.TrainingSummaryView, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var rows []listTrainingSummaryByDateRow
	if _, err = q.Select(&rows, `
	SELECT
		t.id AS training_id,
		t.started_at,
		t.ended_at,
		COUNT(DISTINCT te.id) AS exercise_count,
		COUNT(ts.id) AS set_count
	FROM trainings t
	JOIN training_exercises te ON te.training_id = t.id
	JOIN training_sets ts ON ts.training_exercise_id = te.id
	WHERE t.user_id = ?
	-- この日のトレーニングを取得する
	AND DATE(CONVERT_TZ(t.started_at, '+00:00', '+09:00')) = ?
	GROUP BY t.id
	ORDER BY t.started_at ASC
	`, userIDBytes, date.Format("2006-01-02")); err != nil {
		return nil, err
	}
	var result []*trainingdomain.TrainingSummaryView
	for _, row := range rows {
		summary, err := toTrainingSummaryViewFromRow(row)
		if err != nil {
			return nil, err
		}
		result = append(result, summary)
	}
	return result, nil
}

type getTrainingCountInAWeekRow struct {
	TrainingCount int `db:"training_count"`
}

func (s *trainingQueryService) GetTrainingCountInAWeek(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (valueobject.NonNegativeInt, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return valueobject.NonNegativeInt{}, err
	}
	var row getTrainingCountInAWeekRow
	if err := q.SelectOne(&row, `SELECT COUNT(DISTINCT t.id) AS training_count FROM trainings t WHERE t.user_id = ? AND DATE(CONVERT_TZ(t.started_at, '+00:00', '+09:00')) BETWEEN ? AND ?`, userIDBytes, weekStart.Format("2006-01-02"), weekStart.AddDate(0, 0, 6).Format("2006-01-02")); err != nil {
		return valueobject.NonNegativeInt{}, err
	}
	trainingCount, err := valueobject.NewNonNegativeInt(row.TrainingCount)
	if err != nil {
		return valueobject.NonNegativeInt{}, err
	}
	return *trainingCount, nil
}

func toTrainingSummaryViewFromRow(row listTrainingSummaryByDateRow) (*trainingdomain.TrainingSummaryView, error) {
	trainingID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.TrainingID](row.TrainingID)
	if err != nil {
		return nil, err
	}
	exerciseCount, err := valueobject.NewNonNegativeInt(row.ExerciseCount)
	if err != nil {
		return nil, err
	}
	setCount, err := valueobject.NewNonNegativeInt(row.SetCount)
	if err != nil {
		return nil, err
	}
	return &trainingdomain.TrainingSummaryView{
		TrainingID:    *trainingID,
		StartedAt:     row.StartedAt,
		EndedAt:       row.EndedAt,
		ExerciseCount: *exerciseCount,
		SetCount:      *setCount,
	}, nil
}
