package trainingdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type RoutineRepository interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.RoutineID, userID valueobject.UserID) (*Routine, error)
	FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*Routine, pagination.OffsetPaginator, error)
	Save(ctx context.Context, routine *Routine) error
	DeleteByID(ctx context.Context, id valueobject.RoutineID) error
}
