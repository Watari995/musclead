package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type ListWeightSummaryByDateInput struct {
	UserID valueobject.UserID
	Date   time.Time
}

type ListWeightSummaryByDateOutput struct {
	WeightSummaries []*weightdomain.WeightSummaryView
}

type ListWeightSummaryByDate struct {
	weightQuery weightdomain.WeightQueryService
}

func NewListWeightSummaryByDate(weightQuery weightdomain.WeightQueryService) *ListWeightSummaryByDate {
	return &ListWeightSummaryByDate{weightQuery: weightQuery}
}

func (uc *ListWeightSummaryByDate) Execute(ctx context.Context, input ListWeightSummaryByDateInput) (*ListWeightSummaryByDateOutput, error) {
	weightSummaries, err := uc.weightQuery.ListSummaryByDate(ctx, input.UserID, input.Date)
	if err != nil {
		return nil, err
	}

	return &ListWeightSummaryByDateOutput{WeightSummaries: weightSummaries}, nil
}
