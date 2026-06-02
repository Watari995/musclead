package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type RecordTrainingInput struct {
	UserID       valueobject.UserID
	TrainingSpec trainingdomain.TrainingSpec
}

type RecordTrainingOutput struct {
	TrainingID valueobject.TrainingID
}

type RecordTraining struct {
	trainingRepo trainingdomain.TrainingRepository
	txManager    dbtx.TransactionManager
}

func (uc *RecordTraining) Execute(ctx context.Context, input RecordTrainingInput) (*RecordTrainingOutput, error) {
	training := trainingdomain.CreateTraining(input.TrainingSpec, input.UserID)

	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		if err := uc.trainingRepo.Save(txCtx, training); err != nil {
			return myerror.NewInternalError().Wrap(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &RecordTrainingOutput{TrainingID: training.ID()}, nil
}

func NewRecordTraining(trainingRepo trainingdomain.TrainingRepository, txManager dbtx.TransactionManager) *RecordTraining {
	return &RecordTraining{
		trainingRepo: trainingRepo,
		txManager:    txManager,
	}
}
