package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type DeleteExerciseByIDInput struct {
	ID     valueobject.ExerciseID
	UserID valueobject.UserID
}

type DeleteExerciseByID struct {
	exerciseRepo trainingdomain.ExerciseRepository
}

func (uc *DeleteExerciseByID) Execute(ctx context.Context, input DeleteExerciseByIDInput) error {
	ex, err := uc.exerciseRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if ex == nil {
		return myerror.NewExerciseNotFoundError()
	}

	if err := uc.exerciseRepo.DeleteByID(ctx, ex.ID()); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.Training.ExerciseUsedInTrainingError) {
			return err // repository が myerror を返してくるのでそのまま素通し
		}
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewDeleteExerciseByID(exerciseRepo trainingdomain.ExerciseRepository) *DeleteExerciseByID {
	return &DeleteExerciseByID{exerciseRepo: exerciseRepo}
}
