package trainingusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListTrainingDatesByMonthInput struct {
	UserID valueobject.UserID
	Year   int
	Month  int
}

type ListTrainingDatesByMonthOutput struct {
	Dates []time.Time
}

type ListTrainingDatesByMonth struct {
	trainingQuery trainingdomain.TrainingQueryService
}

func NewListTrainingDatesByMonth(trainingQuery trainingdomain.TrainingQueryService) *ListTrainingDatesByMonth {
	return &ListTrainingDatesByMonth{trainingQuery: trainingQuery}
}

func (uc *ListTrainingDatesByMonth) Execute(ctx context.Context, input ListTrainingDatesByMonthInput) (*ListTrainingDatesByMonthOutput, error) {
	dates, err := uc.trainingQuery.ListTrainingDatesByMonth(ctx, input.UserID, input.Year, input.Month)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListTrainingDatesByMonthOutput{Dates: dates}, nil
}
