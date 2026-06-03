package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateRoutineInput struct {
	UserID      valueobject.UserID
	RoutineSpec trainingdomain.RoutineSpec
}

type CreateRoutineOutput struct {
	ID valueobject.RoutineID
}

type CreateRoutine struct {
	routineRepo trainingdomain.RoutineRepository
}

func (uc *CreateRoutine) Execute(ctx context.Context, input CreateRoutineInput) (*CreateRoutineOutput, error) {
	routine := trainingdomain.CreateRoutine(input.RoutineSpec, input.UserID)
	if err := uc.routineRepo.Save(ctx, routine); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.Training.RoutineNameAlreadyExistsError) {
			return nil, err
		}
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &CreateRoutineOutput{ID: routine.ID()}, nil
}

func NewCreateRoutine(routineRepo trainingdomain.RoutineRepository) *CreateRoutine {
	return &CreateRoutine{routineRepo: routineRepo}
}
