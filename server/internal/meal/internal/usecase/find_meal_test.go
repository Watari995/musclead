package mealusecase_test

import (
	"context"
	"errors"
	"testing"

	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindMealByID_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryId[valueobject.UserID]()
	meal := newDummyMeal(t, userID)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(meal, nil)

	uc := mealusecase.NewFindMealByID(repo)
	output, err := uc.Execute(context.Background(), mealusecase.FindMealByIDInput{
		MealID: meal.ID(),
		UserID: userID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, meal.ID().Value(), output.Meal.ID().Value())
}

func TestFindMealByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil)

	uc := mealusecase.NewFindMealByID(repo)
	output, err := uc.Execute(context.Background(), mealusecase.FindMealByIDInput{
		MealID: valueobject.NewPrimaryId[valueobject.MealID](),
		UserID: valueobject.NewPrimaryId[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Meal.NotFoundError))
}

func TestFindMealByID_OwnerMismatch(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	ownerID := valueobject.NewPrimaryId[valueobject.UserID]()
	otherID := valueobject.NewPrimaryId[valueobject.UserID]()
	meal := newDummyMeal(t, ownerID)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(meal, nil)

	uc := mealusecase.NewFindMealByID(repo)
	output, err := uc.Execute(context.Background(), mealusecase.FindMealByIDInput{
		MealID: meal.ID(),
		UserID: otherID,
	})

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.General.PermissionError))
}

func TestFindMealByID_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("db down"))

	uc := mealusecase.NewFindMealByID(repo)
	output, err := uc.Execute(context.Background(), mealusecase.FindMealByIDInput{
		MealID: valueobject.NewPrimaryId[valueobject.MealID](),
		UserID: valueobject.NewPrimaryId[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.Nil(t, output)
}
