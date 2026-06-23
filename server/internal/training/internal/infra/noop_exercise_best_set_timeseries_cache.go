package traininginfra

import (
	"context"
	"time"

	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// NoOpExerciseBestSetTimeseriesCache はキャッシュを使わない実装。
// ローカル開発・テスト用。training.go で Redis が不要な場合に差し替える。
type NoOpExerciseBestSetTimeseriesCache struct{}

func NewNoOpExerciseBestSetTimeseriesCache() trainingdomain.ExerciseBestSetTimeseriesCache {
	return &NoOpExerciseBestSetTimeseriesCache{}
}

func (c *NoOpExerciseBestSetTimeseriesCache) FindByPeriod(_ context.Context, _ valueobject.UserID, _ valueobject.ExerciseID, _ time.Time, _ time.Time) ([]*trainingdomain.BestSetView, bool, error) {
	return nil, false, nil
}

func (c *NoOpExerciseBestSetTimeseriesCache) Save(_ context.Context, _ valueobject.UserID, _ *trainingdomain.BestSetView) error {
	return nil
}

func (c *NoOpExerciseBestSetTimeseriesCache) Evict(_ context.Context, _ valueobject.UserID, _ valueobject.ExerciseID) error {
	return nil
}
