package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindBestSetsByExerciseIDsInput struct {
	UserID      valueobject.UserID
	ExerciseIDs []valueobject.ExerciseID
}

type FindBestSetsByExerciseIDsOutput struct {
	BestSets []*trainingdomain.BestSetView
}

type FindBestSetsByExerciseIDs struct {
	exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService
}

func (uc *FindBestSetsByExerciseIDs) Execute(ctx context.Context, input FindBestSetsByExerciseIDsInput) (*FindBestSetsByExerciseIDsOutput, error) {
	// 0件チェックはinfraで行う
	bestSets, err := uc.exerciseRecordQueryService.FindBestSetsByExerciseIDs(ctx, input.UserID, input.ExerciseIDs)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &FindBestSetsByExerciseIDsOutput{BestSets: bestSets}, nil
}

func NewFindBestSetsByExerciseIDs(exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService) *FindBestSetsByExerciseIDs {
	return &FindBestSetsByExerciseIDs{exerciseRecordQueryService: exerciseRecordQueryService}
}
