package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type DeleteRoutineByIDInput struct {
	ID     valueobject.RoutineID
	UserID valueobject.UserID
}

type DeleteRoutineByID struct {
	routineRepo trainingdomain.RoutineRepository
}

func (uc *DeleteRoutineByID) Execute(ctx context.Context, input DeleteRoutineByIDInput) error {
	routine, err := uc.routineRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if routine == nil {
		return myerror.NewRoutineNotFoundError()
	}
	if err := uc.routineRepo.DeleteByID(ctx, input.ID); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewDeleteRoutineByID(routineRepo trainingdomain.RoutineRepository) *DeleteRoutineByID {
	return &DeleteRoutineByID{routineRepo: routineRepo}
}
