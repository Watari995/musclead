package trainingusecase

import (
	"context"
	"log/slog"
	"time"

	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetExerciseBestSetTimeseriesInput struct {
	UserID     valueobject.UserID
	ExerciseID valueobject.ExerciseID
	From       time.Time
	To         time.Time
}

type GetExerciseBestSetTimeseriesOutput struct {
	BestSets []*trainingdomain.BestSetView
}

type GetExerciseBestSetTimeseries struct {
	exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService
	cache                      trainingdomain.ExerciseBestSetTimeseriesCache
}

func (uc *GetExerciseBestSetTimeseries) Execute(ctx context.Context, input GetExerciseBestSetTimeseriesInput) (*GetExerciseBestSetTimeseriesOutput, error) {
	// TODO: weight の GetWeightTimeseries.Execute と同じパターンで実装する。
	//
	// 1. cache.FindByPeriod を呼ぶ。hit=true なら即返す。
	// 2. キャッシュミス（hit=false）または error の場合は exerciseRecordQueryService.FindBestSetTimeseriesByExerciseID で DB から取得。
	// 3. DB 取得後、goroutine で populateBestSetCache を呼んでキャッシュを非同期 populate する。
	//    （caller の ctx はレスポンス後にキャンセルされるため context.Background() を使う）
	// 4. output を返す。
	panic("not implemented")
}

func populateBestSetCache(cache trainingdomain.ExerciseBestSetTimeseriesCache, bestSets []*trainingdomain.BestSetView) {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("best set cache populate panicked", "panic", r)
		}
	}()
	bgCtx := context.Background()
	for _, b := range bestSets {
		if err := cache.Save(bgCtx, b); err != nil {
			slog.Warn("best set cache populate failed", "err", err, "trainingID", b.TrainingID.Value())
		}
	}
}

func NewGetExerciseBestSetTimeseries(
	exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService,
	cache trainingdomain.ExerciseBestSetTimeseriesCache,
) *GetExerciseBestSetTimeseries {
	return &GetExerciseBestSetTimeseries{
		exerciseRecordQueryService: exerciseRecordQueryService,
		cache:                      cache,
	}
}

