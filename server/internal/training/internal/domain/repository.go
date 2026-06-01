package trainingdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type TrainingRepository interface {
	Save(ctx context.Context, training *Training) error
	FindByID(ctx context.Context, id valueobject.TrainingID) (*Training, error)
	FindAllByUserIDWithOffsetPagination(
		ctx context.Context,
		userID valueobject.UserID,
		limit int,
		offset int,
	) ([]*Training, pagination.OffsetPaginator, error)
	DeleteByID(ctx context.Context, id valueobject.TrainingID) error
}
