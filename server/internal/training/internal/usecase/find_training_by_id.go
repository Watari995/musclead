package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindTrainingByIDInput struct {
	TrainingID valueobject.TrainingID
	UserID     valueobject.UserID
}

type FindTrainingByIDOutput struct {
	Training *trainingdomain.Training
}

type FindTrainingByID struct {
	trainingRepo trainingdomain.TrainingRepository
}

func (uc *FindTrainingByID) Execute(ctx context.Context, input FindTrainingByIDInput) (*FindTrainingByIDOutput, error) {
	training, err := uc.trainingRepo.FindByID(ctx, input.TrainingID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if training == nil {
		return nil, myerror.NewTrainingNotFoundError()
	}
	if training.UserID() != input.UserID {
		return nil, myerror.NewPermissionError().SetMessage("training does not belong to the user")
	}

	return &FindTrainingByIDOutput{
		Training: training,
	}, nil
}

func NewFindTrainingByID(trainingRepo trainingdomain.TrainingRepository) *FindTrainingByID {
	return &FindTrainingByID{trainingRepo: trainingRepo}
}
