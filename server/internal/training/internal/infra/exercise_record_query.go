package traininginfra

import (
	"context"
	"fmt"
	"time"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlquery"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

// ExerciseRecordQueryService の実装。
// dbtx.Querier + gorp で生 SQL を実行し、BestSetView を組み立てる。
// 実装パターン(struct / New 関数 / row struct / sqlconv 変換)は infra/routine_query.go を参照。
//
//   - exerciseRecordQueryService struct{ dbmap *gorp.DbMap } と NewExerciseRecordQueryService
//   - FindBestSet 用 SQL:
//       SELECT ts.weight_kg, ts.reps, t.started_at, t.id
//       FROM training_sets ts
//       JOIN training_exercises te ON ts.training_exercise_id = te.id
//       JOIN trainings t           ON te.training_id = t.id
//       WHERE t.user_id = ? AND te.exercise_id = ?
//       ORDER BY ts.weight_kg DESC, ts.reps DESC
//       LIMIT 1
//   - 0 件のときは (nil, nil) を返す

type exerciseRecordQueryService struct {
	dbmap *gorp.DbMap
}

func NewExerciseRecordQueryService(dbmap *gorp.DbMap) trainingdomain.ExerciseRecordQueryService {
	return &exerciseRecordQueryService{dbmap: dbmap}
}

// window関数を使用して、各種目の最重量を取得する
func buildFindBestSetsByExerciseIDsSQL(exerciseIDs [][]byte) (string, []any) {
	placeholders, args := sqlquery.InPlaceholders(exerciseIDs)
	return fmt.Sprintf(`
	SELECT weight_kg, reps, started_at, id, exercise_id
	FROM (
		SELECT
			ts.weight_kg AS weight_kg,
			ts.reps      AS reps,
			t.started_at AS started_at,
			t.id          AS id,
			te.exercise_id AS exercise_id,
			ROW_NUMBER() OVER (
				PARTITION BY te.exercise_id
				ORDER BY ts.weight_kg DESC, ts.reps DESC
			) AS rn
		FROM training_sets ts
		JOIN training_exercises te ON ts.training_exercise_id = te.id
		JOIN trainings t ON te.training_id = t.id
		WHERE t.user_id = ? AND te.exercise_id IN (%s)
	) ranked
	WHERE ranked.rn = 1
	`, placeholders), args
}

type bestSetsByExerciseIDsRow struct {
	WeightKg   string    `db:"weight_kg"`
	Reps       int32     `db:"reps"`
	StartedAt  time.Time `db:"started_at"`
	TrainingID []byte    `db:"id"`
	ExerciseID []byte    `db:"exercise_id"`
}

func (s *exerciseRecordQueryService) FindBestSetsByExerciseIDs(ctx context.Context, userID valueobject.UserID, exerciseIDs []valueobject.ExerciseID) ([]*trainingdomain.BestSetView, error) {
	// 0 件のときはエラーになるため、先にチェックして nil を返す
	if len(exerciseIDs) == 0 {
		return nil, nil
	}
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	exerciseIDBytes := make([][]byte, 0, len(exerciseIDs))
	for _, exerciseID := range exerciseIDs {
		bytes, err := exerciseID.Bytes()
		if err != nil {
			return nil, err
		}
		exerciseIDBytes = append(exerciseIDBytes, bytes)
	}
	var rows []bestSetsByExerciseIDsRow
	sqlStr, inArgs := buildFindBestSetsByExerciseIDsSQL(exerciseIDBytes)
	// userID を先頭に追加する
	args := append([]any{userIDBytes}, inArgs...)
	if _, err = q.Select(&rows, sqlStr, args...); err != nil {
		return nil, err
	}
	result := make([]*trainingdomain.BestSetView, 0, len(rows))
	for _, row := range rows {
		bestSet, err := toBestSetViewFromRow(row)
		if err != nil {
			return nil, err
		}
		result = append(result, bestSet)
	}
	return result, nil
}

func buildFindBestSetsByExerciseIDSQL() string {
	return `
	SELECT weight_kg, reps, started_at, id, exercise_id
	FROM (
		SELECT
			ts.weight_kg AS weight_kg,
			ts.reps AS reps,
			t.started_at AS started_at,
			t.id AS id,
			te.exercise_id AS exercise_id,
			ROW_NUMBER() OVER (
				PARTITION BY t.id
				ORDER BY ts.weight_kg DESC, ts.reps DESC
			) AS rn
			FROM training_sets ts
			JOIN training_exercises te ON ts.training_exercise_id = te.id
			JOIN trainings t ON te.training_id = t.id
			WHERE t.user_id = ? AND te.exercise_id = ?
			AND t.started_at BETWEEN ? AND ?
	) ranked
	WHERE ranked.rn = 1
	ORDER BY started_at ASC
	`
}

func (s *exerciseRecordQueryService) FindBestSetTimeseriesByExerciseID(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID, from, to time.Time) ([]*trainingdomain.BestSetView, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	exerciseIDBytes, err := exerciseID.Bytes()
	if err != nil {
		return nil, err
	}

	sqlStr := buildFindBestSetsByExerciseIDSQL()

	var rows []bestSetsByExerciseIDsRow
	if _, err = q.Select(&rows, sqlStr, userIDBytes, exerciseIDBytes, from, to); err != nil {
		return nil, err
	}

	result := make([]*trainingdomain.BestSetView, 0, len(rows))
	for _, row := range rows {
		bestSet, err := toBestSetViewFromRow(row)
		if err != nil {
			return nil, err
		}
		result = append(result, bestSet)
	}
	return result, nil
}

func toBestSetViewFromRow(row bestSetsByExerciseIDsRow) (*trainingdomain.BestSetView, error) {
	weightKg, err := valueobject.NewNonNegativeDecimalFromString(row.WeightKg)
	if err != nil {
		return nil, err
	}
	reps, err := valueobject.NewNonNegativeInt(int(row.Reps))
	if err != nil {
		return nil, err
	}
	trainingID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.TrainingID](row.TrainingID)
	if err != nil {
		return nil, err
	}
	exerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.ExerciseID)
	if err != nil {
		return nil, err
	}
	return &trainingdomain.BestSetView{
		WeightKg:    *weightKg,
		Reps:        *reps,
		PerformedAt: row.StartedAt,
		TrainingID:  *trainingID,
		ExerciseID:  *exerciseID,
	}, nil
}

func buildFindLastSessionSetsByExerciseIDsSQL(exerciseIDs [][]byte) (string, []any) {
	placeholders, args := sqlquery.InPlaceholders(exerciseIDs)
	return fmt.Sprintf(`
	SELECT exercise_id, performed_at, weight_kg, reps, set_number
	FROM (
		SELECT
		  te.exercise_id AS exercise_id,
			t.started_at AS performed_at,
			ts.weight_kg AS weight_kg,
			ts.reps AS reps,
			ts.set_number AS set_number,
			DENSE_RANK() OVER (
				PARTITION BY te.exercise_id
				ORDER BY t.started_at DESC -- 直近のセッションを一つ取得する
			) AS rn
		FROM training_sets ts
		JOIN training_exercises te ON ts.training_exercise_id = te.id
		JOIN trainings t ON te.training_id = t.id
		-- 途中でセッションを終了した時などバグりそうなので今日のセッションは除外する
		WHERE t.user_id = ? AND te.exercise_id IN (%s) AND DATE(t.started_at) <> CURDATE()
	) ranked
	WHERE ranked.rn = 1
	ORDER BY exercise_id, set_number
	`, placeholders), args
}

type lastSessionSetsRow struct {
	ExerciseID  []byte    `db:"exercise_id"`
	PerformedAt time.Time `db:"performed_at"`
	WeightKg    string    `db:"weight_kg"`
	Reps        int32     `db:"reps"`
	SetNumber   int32     `db:"set_number"`
}

func (s *exerciseRecordQueryService) FindLastSessionSetsByExerciseIDs(ctx context.Context, userID valueobject.UserID, exerciseIDs []valueobject.ExerciseID) ([]*trainingdomain.LastSessionSetByExerciseView, error) {
	// 0件の時はerrorになるため先にチェックしてnilを返す
	if len(exerciseIDs) == 0 {
		return nil, nil
	}
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	exerciseIDBytes := make([][]byte, 0, len(exerciseIDs))
	for _, e := range exerciseIDs {
		bytes, err := e.Bytes()
		if err != nil {
			return nil, err
		}
		exerciseIDBytes = append(exerciseIDBytes, bytes)
	}
	var rows []lastSessionSetsRow
	sqlStr, inArgs := buildFindLastSessionSetsByExerciseIDsSQL(exerciseIDBytes)
	// userIDを先頭に追加する
	args := append([]any{userIDBytes}, inArgs...)
	if _, err = q.Select(&rows, sqlStr, args...); err != nil {
		return nil, err
	}

	byExerciseID := map[string]*trainingdomain.LastSessionSetByExerciseView{}
	for _, r := range rows {
		lastSessionSet, err := toLastSessionSetByExerciseViewFromRow(r)
		if err != nil {
			return nil, err
		}
		if v, ok := byExerciseID[lastSessionSet.ExerciseID.String()]; !ok {
			byExerciseID[lastSessionSet.ExerciseID.String()] = lastSessionSet
		} else {
			v.Sets = append(v.Sets, lastSessionSet.Sets[0])
		}
	}
	result := make([]*trainingdomain.LastSessionSetByExerciseView, 0, len(byExerciseID))
	for _, v := range byExerciseID {
		result = append(result, v)
	}
	return result, nil
}

func toLastSessionSetByExerciseViewFromRow(row lastSessionSetsRow) (*trainingdomain.LastSessionSetByExerciseView, error) {
	exerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.ExerciseID)
	if err != nil {
		return nil, err
	}
	setNumber, err := valueobject.NewNonNegativeInt(int(row.SetNumber))
	if err != nil {
		return nil, err
	}
	weightKg, err := valueobject.NewNonNegativeDecimalFromString(row.WeightKg)
	if err != nil {
		return nil, err
	}
	reps, err := valueobject.NewNonNegativeInt(int(row.Reps))
	if err != nil {
		return nil, err
	}
	return &trainingdomain.LastSessionSetByExerciseView{
		ExerciseID:  *exerciseID,
		PerformedAt: row.PerformedAt,
		Sets: []*trainingdomain.LastSessionSetView{
			{
				SetNumber: *setNumber,
				WeightKg:  *weightKg,
				Reps:      *reps,
			},
		},
	}, nil
}
