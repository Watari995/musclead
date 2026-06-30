package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type GetWeightChangeInAWeek struct {
	weightQuery weightdomain.WeightQueryService
}

func NewGetWeightChangeInAWeek(weightQuery weightdomain.WeightQueryService) *GetWeightChangeInAWeek {
	return &GetWeightChangeInAWeek{weightQuery: weightQuery}
}

func (uc *GetWeightChangeInAWeek) Execute(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (*valueobject.WeightChangeKg, error) {
	change, err := uc.weightQuery.GetWeightChangeInAWeek(ctx, userID, weekStart)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return change, nil
}
