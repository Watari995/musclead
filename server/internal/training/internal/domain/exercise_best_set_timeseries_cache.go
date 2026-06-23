package trainingdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// ExerciseBestSetTimeseriesCache は種目ベストセット時系列のキャッシュインターフェース。
// weight の WeightTimeseriesCache と同じ ZSet + Hash 2層構造で Redis に格納する。
//
// キー空間:
//   - exercises:best-set-timeseries:{userID}:{exerciseID}:idx  (ZSet, score=started_at.Unix())
//   - exercises:best-set-timeseries:{userID}:{exerciseID}:data (Hash, field=trainingID, value=JSON)
type ExerciseBestSetTimeseriesCache interface {
	// FindByPeriod は期間内のベストセット時系列を古い順で返す。キャッシュミス時は hit=false。
	FindByPeriod(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID, from, to time.Time) (bestSets []*BestSetView, hit bool, err error)
	// Save は1件のデータポイントをキャッシュに追記し TTL をリセットする。
	// GetExerciseBestSetTimeseries のキャッシュ populate で goroutine から呼ぶ。
	// BestSetView に UserID がないため userID を別引数で受け取る。
	Save(ctx context.Context, userID valueobject.UserID, bestSet *BestSetView) error
	// Evict は指定種目のキャッシュ(idx + data)を全削除する。
	// RecordTraining / UpdateTraining / DeleteTrainingByID の write 後に呼ぶ。
	Evict(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID) error
}
