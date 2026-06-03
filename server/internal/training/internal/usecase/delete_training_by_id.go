package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type DeleteTrainingByIDInput struct {
	TrainingID valueobject.TrainingID
	UserID     valueobject.UserID
}

type DeleteTrainingByID struct {
	trainingRepo trainingdomain.TrainingRepository
}

func (uc *DeleteTrainingByID) Execute(ctx context.Context, input DeleteTrainingByIDInput) error {
	training, err := uc.trainingRepo.FindByID(ctx, input.TrainingID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if training == nil {
		return myerror.NewTrainingNotFoundError().SetMessage("training not found")
	}
	if training.UserID() != input.UserID {
		return myerror.NewPermissionError().SetMessage("training does not belong to the user")
	}

	if err := uc.trainingRepo.DeleteByID(ctx, input.TrainingID); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewDeleteTrainingByID(trainingRepo trainingdomain.TrainingRepository) *DeleteTrainingByID {
	return &DeleteTrainingByID{trainingRepo: trainingRepo}
}
