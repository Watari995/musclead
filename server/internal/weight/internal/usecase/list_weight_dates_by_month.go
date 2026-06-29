package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type ListWeightDatesByMonthInput struct {
	UserID valueobject.UserID
	Year   int
	Month  int
}

type ListWeightDatesByMonthOutput struct {
	Dates []time.Time
}

type ListWeightDatesByMonth struct {
	weightQuery weightdomain.WeightQueryService
}

func NewListWeightDatesByMonth(weightQuery weightdomain.WeightQueryService) *ListWeightDatesByMonth {
	return &ListWeightDatesByMonth{weightQuery: weightQuery}
}

func (uc *ListWeightDatesByMonth) Execute(ctx context.Context, input ListWeightDatesByMonthInput) (*ListWeightDatesByMonthOutput, error) {
	dates, err := uc.weightQuery.ListWeightDatesByMonth(ctx, input.UserID, input.Year, input.Month)
	if err != nil {
		return nil, err
	}

	return &ListWeightDatesByMonthOutput{Dates: dates}, nil
}
