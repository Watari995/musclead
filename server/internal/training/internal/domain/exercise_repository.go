package trainingdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ExerciseRepository interface {
	FindByID(ctx context.Context, id valueobject.ExerciseID) (*Exercise, error)
	FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*Exercise, pagination.OffsetPaginator, error)
	ExistsByName(ctx context.Context, userID valueobject.UserID, name valueobject.String50) (bool, error)
	Save(ctx context.Context, exercise *Exercise) error
	DeleteByID(ctx context.Context, id valueobject.ExerciseID) error
}
