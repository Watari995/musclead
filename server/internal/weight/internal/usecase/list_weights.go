package weightusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type ListWeightsInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListWeightsOutput struct {
	Weights    []*weightdomain.Weight
	Pagination pagination.OffsetPaginator
}

type ListWeights struct {
	weightRepo weightdomain.WeightRepository
}

func (uc *ListWeights) Execute(ctx context.Context, input ListWeightsInput) (*ListWeightsOutput, error) {
	weights, paginator, err := uc.weightRepo.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListWeightsOutput{Weights: weights, Pagination: paginator}, nil
}

func NewListWeights(weightRepo weightdomain.WeightRepository) *ListWeights {
	return &ListWeights{weightRepo: weightRepo}
}
