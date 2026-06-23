package trainingusecase

import (
	"context"
	"log/slog"

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
	// bestSetCache はトレーニング削除後に種目キャッシュを evict するために使う。
	bestSetCache trainingdomain.ExerciseBestSetTimeseriesCache
}

func (uc *DeleteTrainingByID) Execute(ctx context.Context, input DeleteTrainingByIDInput) error {
	training, err := uc.trainingRepo.FindByIDAndUserID(ctx, input.TrainingID, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if training == nil {
		return myerror.NewTrainingNotFoundError().SetMessage("training not found")
	}

	if err := uc.trainingRepo.DeleteByID(ctx, input.TrainingID); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	// トレーニング削除後に種目キャッシュをevictする
	for _, e := range training.Exercises() {
		if err := uc.bestSetCache.Evict(ctx, input.UserID, e.ExerciseID()); err != nil {
			slog.Warn("best set cache evict failed", "err", err)
		}
	}

	return nil
}

func NewDeleteTrainingByID(trainingRepo trainingdomain.TrainingRepository, bestSetCache trainingdomain.ExerciseBestSetTimeseriesCache) *DeleteTrainingByID {
	return &DeleteTrainingByID{trainingRepo: trainingRepo, bestSetCache: bestSetCache}
}
