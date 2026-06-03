package mealusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func recordInput(userID valueobject.UserID) mealusecase.RecordMealInput {
	mealType, _ := valueobject.NewString20("lunch")
	calories, _ := valueobject.NewNonNegativeInt(600)
	return mealusecase.RecordMealInput{
		UserID:   userID,
		EatenAt:  time.Now(),
		MealType: *mealType,
		Calories: *calories,
	}
}

func TestRecordMeal_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()

	repo.On("Save", mock.Anything, mock.AnythingOfType("*mealdomain.Meal")).Return(nil)

	uc := mealusecase.NewRecordMeal(repo, fakeTxManager{})
	output, err := uc.Execute(context.Background(), recordInput(userID))

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.MealID.Value())
	repo.AssertExpectations(t)
}

func TestRecordMeal_SaveError(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()

	repo.On("Save", mock.Anything, mock.AnythingOfType("*mealdomain.Meal")).Return(errors.New("db down"))

	uc := mealusecase.NewRecordMeal(repo, fakeTxManager{})
	output, err := uc.Execute(context.Background(), recordInput(userID))

	assert.Error(t, err)
	assert.Nil(t, output)
}
