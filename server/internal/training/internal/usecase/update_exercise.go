package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateExerciseInput struct {
	ID     valueobject.ExerciseID
	UserID valueobject.UserID
	Name   valueobject.String50
}

type UpdateExerciseOutput struct {
	ID valueobject.ExerciseID
}

type UpdateExercise struct {
	exerciseRepo trainingdomain.ExerciseRepository
}

func (uc *UpdateExercise) Execute(ctx context.Context, input UpdateExerciseInput) (*UpdateExerciseOutput, error) {
	found, err := uc.exerciseRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if found == nil {
		return nil, myerror.NewExerciseNotFoundError()
	}
	if found.UserID() != input.UserID {
		return nil, myerror.NewPermissionError().SetMessage("exercise does not belong to the user")
	}

	found.SetName(input.Name)
	if err := uc.exerciseRepo.Save(ctx, found); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.Training.ExerciseNameAlreadyExistsError) {
			return nil, err
		}
		return nil, myerror.NewInternalError().Wrap(err)
	}

	return &UpdateExerciseOutput{ID: found.ID()}, nil
}

func NewUpdateExercise(exerciseRepo trainingdomain.ExerciseRepository) *UpdateExercise {
	return &UpdateExercise{exerciseRepo: exerciseRepo}
}
