package mealdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type MealTemplateRepository interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.MealTemplateID, userID valueobject.UserID) (*MealTemplate, error)
	FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*MealTemplate, error)
	FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*MealTemplate, pagination.OffsetPaginator, error)
	NextDisplayOrder(ctx context.Context, userID valueobject.UserID) (int, error)
	Save(ctx context.Context, mealTemplate *MealTemplate) error
	DeleteByID(ctx context.Context, id valueobject.MealTemplateID) error
}
