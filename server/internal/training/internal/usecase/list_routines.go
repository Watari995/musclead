package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListRoutineInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListRoutineOutput struct {
	Routines   []*trainingdomain.RoutineView
	Pagination pagination.OffsetPaginator
}

type ListRoutine struct {
	routineQueryService trainingdomain.RoutineQueryService
}

func (uc *ListRoutine) Execute(ctx context.Context, input ListRoutineInput) (*ListRoutineOutput, error) {
	routines, paginator, err := uc.routineQueryService.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListRoutineOutput{Routines: routines, Pagination: paginator}, nil
}

func NewListRoutine(routineQueryService trainingdomain.RoutineQueryService) *ListRoutine {
	return &ListRoutine{routineQueryService: routineQueryService}
}
