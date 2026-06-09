package weightusecase

import (
	"context"
	"log/slog"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type RecordWeightInput struct {
	UserID     valueobject.UserID
	WeightSpec weightdomain.WeightSpec
}

type RecordWeightOutput struct {
	WeightID valueobject.WeightID
}

type RecordWeight struct {
	weightRepo  weightdomain.WeightRepository
	weightCache weightdomain.WeightTimeseriesCache
}

func (uc *RecordWeight) Execute(ctx context.Context, input RecordWeightInput) (*RecordWeightOutput, error) {
	weight := weightdomain.CreateWeight(input.UserID, input.WeightSpec)
	if err := uc.weightRepo.Save(ctx, weight); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	// cache save
	if err := uc.weightCache.Save(ctx, weight); err != nil {
		slog.Warn("cache save failed", "err", err, "weight", weight.ID().Value())
	}
	return &RecordWeightOutput{WeightID: weight.ID()}, nil
}

func NewRecordWeight(weightRepo weightdomain.WeightRepository, weightCache weightdomain.WeightTimeseriesCache) *RecordWeight {
	return &RecordWeight{weightRepo: weightRepo, weightCache: weightCache}
}
