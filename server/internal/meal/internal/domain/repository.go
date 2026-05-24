package mealdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type MealRepository interface {
	FindAllByUserIDWithOffsetPagination(ctx context.Context, userId valueobject.UserID, limit int, offset int) ([]*Meal, pagination.OffsetPaginator, error)
	FindByID(ctx context.Context, id valueobject.MealID) (*Meal, error)
	Save(ctx context.Context, meal *Meal) error
	DeleteByID(ctx context.Context, id valueobject.MealID) error
}
