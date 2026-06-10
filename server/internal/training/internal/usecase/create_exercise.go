package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateExerciseInput struct {
	UserID valueobject.UserID
	Name   valueobject.String50
}

type CreateExerciseOutput struct {
	ID valueobject.ExerciseID
}

type CreateExercise struct {
	exerciseRepo trainingdomain.ExerciseRepository
}

func (uc *CreateExercise) Execute(ctx context.Context, input CreateExerciseInput) (*CreateExerciseOutput, error) {
	next, err := uc.exerciseRepo.NextDisplayOrder(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	displayOrder, err := valueobject.NewNonNegativeInt(next)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	exercise := trainingdomain.CreateExercise(input.UserID, input.Name, *displayOrder)
	if err := uc.exerciseRepo.Save(ctx, exercise); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.Training.ExerciseNameAlreadyExistsError) {
			return nil, err
		}
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &CreateExerciseOutput{ID: exercise.ID()}, nil
}

func NewCreateExercise(exerciseRepo trainingdomain.ExerciseRepository) *CreateExercise {
	return &CreateExercise{exerciseRepo: exerciseRepo}
}
