package trainingusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateTrainingInput struct {
	TrainingID valueobject.TrainingID
	UserID     valueobject.UserID
	StartedAt  time.Time
	EndedAt    *time.Time
	Memo       *valueobject.String1000
	Exercises  []trainingdomain.ExerciseSpec
}

type UpdateTrainingOutput struct {
	TrainingID valueobject.TrainingID
}

type UpdateTraining struct {
	trainingRepo trainingdomain.TrainingRepository
	txManager    dbtx.TransactionManager
}

func (uc *UpdateTraining) Execute(ctx context.Context, input UpdateTrainingInput) (*UpdateTrainingOutput, error) {
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
	training.Update(trainingdomain.UpdateParams{
		StartedAt: input.StartedAt,
		EndedAt:   input.EndedAt,
		Memo:      input.Memo,
		Exercises: input.Exercises,
	})
	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		err := uc.trainingRepo.Save(txCtx, training)
		if err != nil {
			return err
		}
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
