package trainingdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type TrainingRepository interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.TrainingID, userID valueobject.UserID) (*Training, error)
	FindAllByUserIDWithOffsetPagination(
		ctx context.Context,
		userID valueobject.UserID,
		limit int,
		offset int,
	) ([]*Training, pagination.OffsetPaginator, error)
	Save(ctx context.Context, training *Training) error
	DeleteByID(ctx context.Context, id valueobject.TrainingID) error
}
