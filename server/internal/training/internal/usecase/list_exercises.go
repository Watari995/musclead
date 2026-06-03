package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListExercisesInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListExercisesOutput struct {
	Exercises  []*trainingdomain.Exercise
	Pagination pagination.OffsetPaginator
}

type ListExercises struct {
	exerciseRepo trainingdomain.ExerciseRepository
}

func (uc *ListExercises) Execute(ctx context.Context, input ListExercisesInput) (*ListExercisesOutput, error) {
	exercises, paginator, err := uc.exerciseRepo.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListExercisesOutput{Exercises: exercises, Pagination: paginator}, nil
}

func NewListExercises(exerciseRepo trainingdomain.ExerciseRepository) *ListExercises {
	return &ListExercises{exerciseRepo: exerciseRepo}
}
