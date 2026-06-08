package weightusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type FindWeightByIDInput struct {
	ID     valueobject.WeightID
	UserID valueobject.UserID
}

type FindWeightByIDOutput struct {
	Weight *weightdomain.Weight
}

type FindWeightByID struct {
	weightRepo weightdomain.WeightRepository
}

func (uc *FindWeightByID) Execute(ctx context.Context, input FindWeightByIDInput) (*FindWeightByIDOutput, error) {
	weight, err := uc.weightRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if weight == nil {
		return nil, myerror.NewWeightNotFoundError()
	}
	return &FindWeightByIDOutput{Weight: weight}, nil
}

func NewFindWeightByID(weightRepo weightdomain.WeightRepository) *FindWeightByID {
	return &FindWeightByID{weightRepo: weightRepo}
}
