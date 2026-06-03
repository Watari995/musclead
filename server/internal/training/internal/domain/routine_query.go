package trainingdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type RoutineExerciseView struct {
	ID           valueobject.RoutineExerciseID
	ExerciseID   valueobject.ExerciseID
	ExerciseName valueobject.String50
	DisplayOrder valueobject.NonNegativeInt
}

type RoutineView struct {
	ID               valueobject.RoutineID
	UserID           valueobject.UserID
	Name             valueobject.String50
	CreatedAt        time.Time
	UpdatedAt        time.Time
	RoutineExercises []RoutineExerciseView
}

type RoutineQueryService interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.RoutineID, userID valueobject.UserID) (*RoutineView, error)
	FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]*RoutineView, pagination.OffsetPaginator, error)
}
