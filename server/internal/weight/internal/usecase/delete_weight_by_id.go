package weightusecase

import (
	"context"
	"log/slog"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type DeleteWeightByIDInput struct {
	ID     valueobject.WeightID
	UserID valueobject.UserID
}

type DeleteWeightByID struct {
	weightRepo  weightdomain.WeightRepository
	weightCache weightdomain.WeightTimeseriesCache
}

func (uc *DeleteWeightByID) Execute(ctx context.Context, input DeleteWeightByIDInput) error {
	weight, err := uc.weightRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if weight == nil {
		return myerror.NewWeightNotFoundError()
	}
	if err := uc.weightRepo.DeleteByID(ctx, input.ID); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if err := uc.weightCache.Delete(ctx, input.UserID, input.ID); err != nil {
		slog.Warn("cache delete failed", "err", err, "weight", input.ID.Value())
	}
	return nil
}

func NewDeleteWeightByID(weightRepo weightdomain.WeightRepository, weightCache weightdomain.WeightTimeseriesCache) *DeleteWeightByID {
	return &DeleteWeightByID{weightRepo: weightRepo, weightCache: weightCache}
}
