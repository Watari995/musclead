package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindExerciseByIDInput struct {
	ID     valueobject.ExerciseID
	UserID valueobject.UserID
}

type FindExerciseByIDOutput struct {
	Exercise *trainingdomain.Exercise
}

type FindExerciseByID struct {
	exerciseRepo trainingdomain.ExerciseRepository
}

func (uc *FindExerciseByID) Execute(ctx context.Context, input FindExerciseByIDInput) (*FindExerciseByIDOutput, error) {
	exercise, err := uc.exerciseRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if exercise == nil {
		return nil, myerror.NewExerciseNotFoundError()
	}
	return &FindExerciseByIDOutput{Exercise: exercise}, nil
}

func NewFindExerciseByID(exerciseRepo trainingdomain.ExerciseRepository) *FindExerciseByID {
	return &FindExerciseByID{exerciseRepo: exerciseRepo}
}
