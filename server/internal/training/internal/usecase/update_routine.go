package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateRoutineInput struct {
	ID          valueobject.RoutineID
	UserID      valueobject.UserID
	RoutineSpec trainingdomain.RoutineSpec
}

type UpdateRoutineOutput struct {
	ID valueobject.RoutineID
}

type UpdateRoutine struct {
	routineRepo trainingdomain.RoutineRepository
}

func (uc *UpdateRoutine) Execute(ctx context.Context, input UpdateRoutineInput) (*UpdateRoutineOutput, error) {
	routine, err := uc.routineRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if routine == nil {
		return nil, myerror.NewRoutineNotFoundError()
	}
	routine.Update(input.RoutineSpec)
	if err := uc.routineRepo.Save(ctx, routine); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.Training.RoutineNameAlreadyExistsError) {
			return nil, err
		}
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpdateRoutineOutput{ID: routine.ID()}, nil
}

func NewUpdateRoutine(routineRepo trainingdomain.RoutineRepository) *UpdateRoutine {
	return &UpdateRoutine{routineRepo: routineRepo}
}
