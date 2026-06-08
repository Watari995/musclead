package weightdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type WeightRepository interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.WeightID, userID valueobject.UserID) (*Weight, error)
	FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*Weight, pagination.OffsetPaginator, error)
	Save(ctx context.Context, weight *Weight) error
	DeleteByID(ctx context.Context, id valueobject.WeightID) error
}
