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

func TestDeleteMealByID_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryId[valueobject.UserID]()
	meal := newDummyMeal(t, userID)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(meal, nil)
	repo.On("DeleteByID", mock.Anything, mock.Anything).Return(nil)

	uc := mealusecase.NewDeleteMealByID(repo)
	err := uc.Execute(context.Background(), mealusecase.DeleteMealByIDInput{
		MealID: meal.ID(),
		UserID: userID,
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteMealByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil)

	uc := mealusecase.NewDeleteMealByID(repo)
	err := uc.Execute(context.Background(), mealusecase.DeleteMealByIDInput{
		MealID: valueobject.NewPrimaryId[valueobject.MealID](),
		UserID: valueobject.NewPrimaryId[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Meal.NotFoundError))
}

func TestDeleteMealByID_OwnerMismatch(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	ownerID := valueobject.NewPrimaryId[valueobject.UserID]()
	otherID := valueobject.NewPrimaryId[valueobject.UserID]()
	meal := newDummyMeal(t, ownerID)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(meal, nil)

	uc := mealusecase.NewDeleteMealByID(repo)
	err := uc.Execute(context.Background(), mealusecase.DeleteMealByIDInput{
		MealID: meal.ID(),
		UserID: otherID,
	})

	assert.Error(t, err)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.General.PermissionError))
}

func TestDeleteMealByID_DeleteError(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryId[valueobject.UserID]()
	meal := newDummyMeal(t, userID)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(meal, nil)
	repo.On("DeleteByID", mock.Anything, mock.Anything).Return(errors.New("db down"))

	uc := mealusecase.NewDeleteMealByID(repo)
	err := uc.Execute(context.Background(), mealusecase.DeleteMealByIDInput{
		MealID: meal.ID(),
		UserID: userID,
	})

	assert.Error(t, err)
}
