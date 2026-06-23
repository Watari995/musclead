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
	// bestSetCache はトレーニング更新後に種目キャッシュを evict するために使う。
	bestSetCache trainingdomain.ExerciseBestSetTimeseriesCache
}

func (uc *UpdateTraining) Execute(ctx context.Context, input UpdateTrainingInput) (*UpdateTrainingOutput, error) {
	training, err := uc.trainingRepo.FindByIDAndUserID(ctx, input.TrainingID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if training == nil {
		return nil, myerror.NewTrainingNotFoundError()
	}

	// TODO: 更新前の種目IDを収集する（preExerciseIDs）。
	//       training.Exercises() から ExerciseID() を取り出す。

	training.Update(input.TrainingSpec)
	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		uc.trainingRepo.Save(txCtx, training)
		return nil
	}); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	// TODO: 更新後の種目IDを収集し（postExerciseIDs）、preExerciseIDs と union して重複排除。
	//       union した全 exerciseID に対して bestSetCache.Evict を呼ぶ。
	//       evict はベストエフォート（失敗しても更新は成功扱い）。slog.Warn でログ。

	return &UpdateTrainingOutput{TrainingID: training.ID()}, nil
}

func NewUpdateTraining(trainingRepo trainingdomain.TrainingRepository, txManager dbtx.TransactionManager, bestSetCache trainingdomain.ExerciseBestSetTimeseriesCache) *UpdateTraining {
	return &UpdateTraining{
		trainingRepo: trainingRepo,
		txManager:    txManager,
		bestSetCache: bestSetCache,
	}
}
