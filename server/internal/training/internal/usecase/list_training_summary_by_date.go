package trainingusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListTrainingSummaryByDateInput struct {
	UserID valueobject.UserID
	Date   time.Time
}

type ListTrainingSummaryByDateOutput struct {
	TrainingSummaries []*trainingdomain.TrainingSummaryView
}

type ListTrainingSummaryByDate struct {
	trainingQuery trainingdomain.TrainingQueryService
}

func NewListTrainingSummaryByDate(trainingQuery trainingdomain.TrainingQueryService) *ListTrainingSummaryByDate {
	return &ListTrainingSummaryByDate{trainingQuery: trainingQuery}
}

func (uc *ListTrainingSummaryByDate) Execute(ctx context.Context, input ListTrainingSummaryByDateInput) (*ListTrainingSummaryByDateOutput, error) {
	summaries, err := uc.trainingQuery.ListSummaryByDate(ctx, input.UserID, input.Date)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListTrainingSummaryByDateOutput{TrainingSummaries: summaries}, nil
}
