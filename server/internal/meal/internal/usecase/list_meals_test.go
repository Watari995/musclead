package mealusecase_test

import (
	"context"
	"errors"
	"testing"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListMeals_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()

	meals := []*mealdomain.Meal{
		newDummyMeal(t, userID),
		newDummyMeal(t, userID),
	}
	pg := pagination.OffsetPaginator{
		CurrentPage: 1, ItemsPerPage: 20, TotalItems: 2, TotalPages: 1,
	}

	repo.On("FindAllByUserIDWithOffsetPagination",
		mock.Anything, mock.Anything, 20, 0,
	).Return(meals, pg, nil)

	uc := mealusecase.NewListMeals(repo)
	output, err := uc.Execute(context.Background(), mealusecase.ListMealsInput{
		UserID: userID,
		Limit:  20,
		Offset: 0,
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Meals, 2)
	assert.Equal(t, 2, output.Pagination.TotalItems)
	repo.AssertExpectations(t)
}

func TestListMeals_Empty(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()

	pg := pagination.OffsetPaginator{
		CurrentPage: 1, ItemsPerPage: 20, TotalItems: 0, TotalPages: 0,
	}
	repo.On("FindAllByUserIDWithOffsetPagination",
		mock.Anything, mock.Anything, 20, 0,
	).Return([]*mealdomain.Meal{}, pg, nil)

	uc := mealusecase.NewListMeals(repo)
	output, err := uc.Execute(context.Background(), mealusecase.ListMealsInput{
		UserID: userID, Limit: 20, Offset: 0,
	})

	assert.NoError(t, err)
	assert.Len(t, output.Meals, 0)
}

func TestListMeals_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	repo.On("FindAllByUserIDWithOffsetPagination",
		mock.Anything, mock.Anything, 20, 0,
	).Return(nil, pagination.OffsetPaginator{}, errors.New("db down"))

	uc := mealusecase.NewListMeals(repo)
	output, err := uc.Execute(context.Background(), mealusecase.ListMealsInput{
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
		Limit:  20, Offset: 0,
	})

	assert.Error(t, err)
	assert.Nil(t, output)
}
