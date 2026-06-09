package weightusecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type GetWeightTimeseriesInput struct {
	UserID valueobject.UserID
	From   time.Time
	To     time.Time
}

type GetWeightTimeseriesOutput struct {
	Weights []*weightdomain.Weight
}

type GetWeightTimeseries struct {
	weightRepo  weightdomain.WeightRepository
	weightCache weightdomain.WeightTimeseriesCache
}

func (uc *GetWeightTimeseries) Execute(ctx context.Context, input GetWeightTimeseriesInput) (*GetWeightTimeseriesOutput, error) {
	// cache miss or errorの時はDBから取得する
	weights, hit, err := uc.weightCache.FindByPeriod(ctx, input.UserID, input.From, input.To)
	if err == nil && hit {
		return &GetWeightTimeseriesOutput{Weights: weights}, nil
	}

	// DB
	weights, err = uc.weightRepo.FindAllByUserIDAndPeriod(ctx, input.UserID, input.From, input.To)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	// cache populate (best effort)
	go populateCache(uc.weightCache, weights)

	return &GetWeightTimeseriesOutput{Weights: weights}, nil
}

func populateCache(cache weightdomain.WeightTimeseriesCache, weights []*weightdomain.Weight) {
	// panic ガード
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("cache populate panicked", "panic", r)
		}
	}()
	// caller ctxはresponse後にキャンセルされるので独立したctxを使う
	bgCtx := context.Background()

	for _, w := range weights {
		if err := cache.Save(bgCtx, w); err != nil {
			slog.Warn("cache populate failed", "err", err, "weight", w.ID().Value())
		}
	}
}

func NewGetWeightTimeseries(weightRepo weightdomain.WeightRepository, weightCache weightdomain.WeightTimeseriesCache) *GetWeightTimeseries {
	return &GetWeightTimeseries{weightRepo: weightRepo, weightCache: weightCache}
}
