package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type DeleteRoutineInput struct {
	ID     valueobject.RoutineID
	UserID valueobject.UserID
}

type DeleteRoutineOutput struct {
	ID valueobject.RoutineID
}

type DeleteRoutine struct {
	routineRepo trainingdomain.RoutineRepository
}

func (uc *DeleteRoutine) Execute(ctx context.Context, input DeleteRoutineInput) (*DeleteRoutineOutput, error) {
	routine, err := uc.routineRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if routine == nil {
		return nil, myerror.NewRoutineNotFoundError()
	}
	if err := uc.routineRepo.DeleteByID(ctx, input.ID); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &DeleteRoutineOutput{ID: routine.ID()}, nil
}

func NewDeleteRoutine(routineRepo trainingdomain.RoutineRepository) *DeleteRoutine {
	return &DeleteRoutine{routineRepo: routineRepo}
}
