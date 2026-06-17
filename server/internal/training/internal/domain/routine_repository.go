package trainingdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type RoutineRepository interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.RoutineID, userID valueobject.UserID) (*Routine, error)
	FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*Routine, error)
	CountByUserID(ctx context.Context, userID valueobject.UserID) (int, error)
	NextDisplayOrder(ctx context.Context, userID valueobject.UserID) (int, error)
	Save(ctx context.Context, routine *Routine) error
	DeleteByID(ctx context.Context, id valueobject.RoutineID) error
}
