package trainingdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ExerciseRepository interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.ExerciseID, userID valueobject.UserID) (*Exercise, error)
	FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*Exercise, error)
	FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*Exercise, pagination.OffsetPaginator, error)
	NextDisplayOrder(ctx context.Context, userID valueobject.UserID) (int, error)
	Save(ctx context.Context, exercise *Exercise) error
	DeleteByID(ctx context.Context, id valueobject.ExerciseID) error
}
