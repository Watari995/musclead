package mealusecase_test

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/mock"
)

// MockMealRepository は mealdomain.MealRepository の偽実装(テスト用)。
type MockMealRepository struct {
	mock.Mock
}

func (m *MockMealRepository) FindAllByUserIDWithOffsetPagination(
	ctx context.Context,
	userId valueobject.UserID,
	limit int,
	offset int,
) ([]*mealdomain.Meal, pagination.OffsetPaginator, error) {
	args := m.Called(ctx, userId, limit, offset)
	meals, _ := args.Get(0).([]*mealdomain.Meal)
	pg, _ := args.Get(1).(pagination.OffsetPaginator)
	return meals, pg, args.Error(2)
}

func (m *MockMealRepository) FindByID(ctx context.Context, id valueobject.MealID) (*mealdomain.Meal, error) {
	args := m.Called(ctx, id)
	meal, _ := args.Get(0).(*mealdomain.Meal)
	return meal, args.Error(1)
}

func (m *MockMealRepository) Save(ctx context.Context, meal *mealdomain.Meal) error {
	args := m.Called(ctx, meal)
	return args.Error(0)
}

func (m *MockMealRepository) DeleteByID(ctx context.Context, id valueobject.MealID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
