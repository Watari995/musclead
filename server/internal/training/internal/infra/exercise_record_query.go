package traininginfra

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
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

const findBestSetSQL = `
SELECT ts.weight_kg, ts.reps, t.started_at, t.id
FROM training_sets ts
JOIN training_exercises te ON ts.training_exercise_id = te.id
JOIN trainings t ON te.training_id = t.id
WHERE t.user_id = ? AND te.exercise_id = ?
-- weight_kg が最大のセットを取得、weight_kg が同じ場合は reps が最大のセットを取得
ORDER BY ts.weight_kg DESC, ts.reps DESC
LIMIT 1
`

type bestSetRow struct {
	WeightKg   string    `db:"weight_kg"`
	Reps       int32     `db:"reps"`
	StartedAt  time.Time `db:"started_at"`
	TrainingID []byte    `db:"id"`
}

func (s *exerciseRecordQueryService) FindBestSet(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID) (*trainingdomain.BestSetView, error) {
	q := dbtx.Querier(ctx, s.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	exerciseIDBytes, err := exerciseID.Bytes()
	if err != nil {
		return nil, err
	}
	var row bestSetRow
	err = q.SelectOne(&row, findBestSetSQL, userIDBytes, exerciseIDBytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 記録なし(nil)は 404 ではなく「記録なし」として正常応答にする想定
		}
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
	trainingID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.TrainingID](row.TrainingID)
	if err != nil {
		return nil, err
	}
	return &trainingdomain.BestSetView{
		WeightKg:    *weightKg,
		Reps:        *reps,
		PerformedAt: row.StartedAt,
		TrainingID:  *trainingID,
	}, nil
}
