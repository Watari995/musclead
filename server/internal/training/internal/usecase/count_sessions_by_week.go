package trainingusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CountSessionsByWeek struct {
	trainingQuery trainingdomain.TrainingQueryService
}

func NewCountSessionsByWeek(trainingQuery trainingdomain.TrainingQueryService) *CountSessionsByWeek {
	return &CountSessionsByWeek{trainingQuery: trainingQuery}
}

func (uc *CountSessionsByWeek) Execute(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (valueobject.NonNegativeInt, error) {
	count, err := uc.trainingQuery.GetTrainingCountInAWeek(ctx, userID, weekStart)
	if err != nil {
		return valueobject.NonNegativeInt{}, myerror.NewInternalError().Wrap(err)
	}
	return count, nil
}
