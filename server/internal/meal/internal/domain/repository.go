package mealdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type MealRepository interface {
	FindByID(ctx context.Context, id valueobject.MealID) (*Meal, error)
	FindByUserID(ctx context.Context, userId valueobject.UserID) ([]*Meal, error)
	Save(ctx context.Context, meal *Meal) error
}
