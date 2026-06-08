package weightusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type UpdateWeightInput struct {
	ID         valueobject.WeightID
	UserID     valueobject.UserID
	WeightSpec weightdomain.WeightSpec
}

type UpdateWeightOutput struct {
	WeightID valueobject.WeightID
}

type UpdateWeight struct {
	weightRepo weightdomain.WeightRepository
}

func (uc *UpdateWeight) Execute(ctx context.Context, input UpdateWeightInput) (*UpdateWeightOutput, error) {
	weight, err := uc.weightRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if weight == nil {
		return nil, myerror.NewWeightNotFoundError()
	}
	weight.Update(input.WeightSpec)
	if err := uc.weightRepo.Save(ctx, weight); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpdateWeightOutput{WeightID: weight.ID()}, nil
}

func NewUpdateWeight(weightRepo weightdomain.WeightRepository) *UpdateWeight {
	return &UpdateWeight{weightRepo: weightRepo}
}
