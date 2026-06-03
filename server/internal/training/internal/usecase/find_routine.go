package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindRoutineByIDInput struct {
	ID     valueobject.RoutineID
	UserID valueobject.UserID
}

type FindRoutineByIDOutput struct {
	Routine *trainingdomain.RoutineView
}

type FindRoutineByID struct {
	routineQueryService trainingdomain.RoutineQueryService
}

func (uc *FindRoutineByID) Execute(ctx context.Context, input FindRoutineByIDInput) (*FindRoutineByIDOutput, error) {
	routineView, err := uc.routineQueryService.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if routineView == nil {
		return nil, myerror.NewRoutineNotFoundError()
	}
	return &FindRoutineByIDOutput{Routine: routineView}, nil
}

func NewFindRoutineByID(routineQueryService trainingdomain.RoutineQueryService) *FindRoutineByID {
	return &FindRoutineByID{routineQueryService: routineQueryService}
}
