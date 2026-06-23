package trainingusecase

import (
	"context"
	"log/slog"

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

	preExerciseIds := make([]valueobject.ExerciseID, 0, len(training.Exercises()))
	for _, e := range training.Exercises() {
		preExerciseIds = append(preExerciseIds, e.ExerciseID())
	}

	training.Update(input.TrainingSpec)

	postExerciseIds := make([]valueobject.ExerciseID, 0, len(training.Exercises()))
	for _, e := range training.Exercises() {
		postExerciseIds = append(postExerciseIds, e.ExerciseID())
	}
	if err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		uc.trainingRepo.Save(txCtx, training)
		return nil
	}); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	uniqueExerciseIdsForEvict := make(map[valueobject.ExerciseID]struct{})
	for _, eid := range preExerciseIds {
		uniqueExerciseIdsForEvict[eid] = struct{}{}
	}
	for _, eid := range postExerciseIds {
		uniqueExerciseIdsForEvict[eid] = struct{}{}
	}
	for eid := range uniqueExerciseIdsForEvict {
		if err := uc.bestSetCache.Evict(ctx, input.UserID, eid); err != nil {
			slog.Warn("best set cache evict failed", "err", err)
		}
	}

	return &UpdateTrainingOutput{TrainingID: training.ID()}, nil
}

func NewUpdateTraining(trainingRepo trainingdomain.TrainingRepository, txManager dbtx.TransactionManager, bestSetCache trainingdomain.ExerciseBestSetTimeseriesCache) *UpdateTraining {
	return &UpdateTraining{
		trainingRepo: trainingRepo,
		txManager:    txManager,
		bestSetCache: bestSetCache,
	}
}
