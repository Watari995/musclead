package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindLastSessionSetsByExerciseIDsInput struct {
	UserID      valueobject.UserID
	ExerciseIDs []valueobject.ExerciseID
}

type FindLastSessionSetsByExerciseIDsOutput struct {
	ExerciseRecords []*trainingdomain.LastSessionSetByExerciseView
}

type FindLastSessionSetsByExerciseIDs struct {
	exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService
}

func (uc *FindLastSessionSetsByExerciseIDs) Execute(ctx context.Context, input FindLastSessionSetsByExerciseIDsInput) (*FindLastSessionSetsByExerciseIDsOutput, error) {
	exerciseRecords, err := uc.exerciseRecordQueryService.FindLastSessionSetsByExerciseIDs(ctx, input.UserID, input.ExerciseIDs)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &FindLastSessionSetsByExerciseIDsOutput{ExerciseRecords: exerciseRecords}, nil
}

func NewFindLastSessionSetsByExerciseIDs(exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService) *FindLastSessionSetsByExerciseIDs {
	return &FindLastSessionSetsByExerciseIDs{exerciseRecordQueryService: exerciseRecordQueryService}
}
