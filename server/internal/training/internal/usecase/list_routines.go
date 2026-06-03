package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListRoutinesInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListRoutinesOutput struct {
	Routines   []*trainingdomain.RoutineView
	Pagination pagination.OffsetPaginator
}

type ListRoutines struct {
	routineQueryService trainingdomain.RoutineQueryService
}

func (uc *ListRoutines) Execute(ctx context.Context, input ListRoutinesInput) (*ListRoutinesOutput, error) {
	routines, paginator, err := uc.routineQueryService.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListRoutinesOutput{Routines: routines, Pagination: paginator}, nil
}

func NewListRoutines(routineQueryService trainingdomain.RoutineQueryService) *ListRoutines {
	return &ListRoutines{routineQueryService: routineQueryService}
}
