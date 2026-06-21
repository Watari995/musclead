package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateTrainingInput struct {
	TrainingID   valueobject.TrainingID
	UserID       valueobject.UserID
	TrainingSpec trainingdomain.TrainingSpec
}

type UpdateTrainingOutput struct {
	TrainingID valueobject.TrainingID
}

type UpdateTraining struct {
	trainingRepo trainingdomain.TrainingRepository
	txManager    dbtx.TransactionManager
}

func (uc *UpdateTraining) Execute(ctx context.Context, input UpdateTrainingInput) (*UpdateTrainingOutput, error) {
	training, err := uc.trainingRepo.FindByIDAndUserID(ctx, input.TrainingID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if training == nil {
		return nil, myerror.NewTrainingNotFoundError()
	}
	training.Update(input.TrainingSpec)
	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		uc.trainingRepo.Save(txCtx, training)
		return nil
	}); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	return &UpdateTrainingOutput{TrainingID: training.ID()}, nil
}

func NewUpdateTraining(trainingRepo trainingdomain.TrainingRepository, txManager dbtx.TransactionManager) *UpdateTraining {
	return &UpdateTraining{
		trainingRepo: trainingRepo,
		txManager:    txManager,
	}
}
