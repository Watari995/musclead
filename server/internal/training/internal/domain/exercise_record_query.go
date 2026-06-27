package trainingdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type BestSetView struct {
	WeightKg    valueobject.NonNegativeDecimal
	Reps        valueobject.NonNegativeInt
	PerformedAt time.Time
	TrainingID  valueobject.TrainingID
	ExerciseID  valueobject.ExerciseID
}

type LastSessionSetByExerciseView struct {
	ExerciseID  valueobject.ExerciseID
	PerformedAt time.Time
	Sets        []*LastSessionSetView
}

type LastSessionSetView struct {
	SetNumber valueobject.NonNegativeInt
	WeightKg  valueobject.NonNegativeDecimal
	Reps      valueobject.NonNegativeInt
}

type ExerciseRecordQueryService interface {
	FindBestSetsByExerciseIDs(ctx context.Context, userID valueobject.UserID, exerciseIDs []valueobject.ExerciseID) ([]*BestSetView, error)
	// FindBestSetTimeseriesByExerciseID は期間内のセッションごとのベストセットを古い順で返す。
	// SQL: PARTITION BY t.id (セッション単位) ORDER BY weight_kg DESC, reps DESC で rn=1 を取り、
	//      t.started_at BETWEEN from AND to で絞り込み、ORDER BY started_at ASC で返す。
	FindBestSetTimeseriesByExerciseID(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID, from, to time.Time) ([]*BestSetView, error)
	// 前回のセッションのセットをexerciseIDごとに返す
	FindLastSessionSetsByExerciseIDs(ctx context.Context, userID valueobject.UserID, exerciseIDs []valueobject.ExerciseID) ([]*LastSessionSetByExerciseView, error)
}
