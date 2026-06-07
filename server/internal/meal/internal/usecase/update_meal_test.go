package mealusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func updateInput(mealID valueobject.MealID, userID valueobject.UserID) mealusecase.UpdateMealInput {
	mealType, _ := valueobject.NewString20("dinner")
	calories, _ := valueobject.NewNonNegativeInt(700)
	return mealusecase.UpdateMealInput{
		MealID:   mealID,
		UserID:   userID,
		EatenAt:  time.Now(),
		MealType: *mealType,
		Calories: *calories,
	}
}

func TestUpdateMeal_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	meal := newDummyMeal(t, userID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(meal, nil)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*mealdomain.Meal")).Return(nil)

	uc := mealusecase.NewUpdateMeal(repo, fakeTxManager{}, fakeStorageClient{})
	output, err := uc.Execute(context.Background(), updateInput(meal.ID(), userID))

	assert.NoError(t, err)
	assert.NotNil(t, output)
	repo.AssertExpectations(t)
}

func TestUpdateMeal_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := mealusecase.NewUpdateMeal(repo, fakeTxManager{}, fakeStorageClient{})
	output, err := uc.Execute(context.Background(), updateInput(
		valueobject.NewPrimaryID[valueobject.MealID](),
		valueobject.NewPrimaryID[valueobject.UserID](),
	))

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Meal.NotFoundError))
}

// 他ユーザーの meal を引こうとした場合: repository が user_id で絞るので nil が返り、
// NotFound として扱われる (他人の meal の存在自体を漏らさない設計)。
func TestUpdateMeal_OtherUserTreatedAsNotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	ownerID := valueobject.NewPrimaryID[valueobject.UserID]()
	otherID := valueobject.NewPrimaryID[valueobject.UserID]()
	meal := newDummyMeal(t, ownerID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := mealusecase.NewUpdateMeal(repo, fakeTxManager{}, fakeStorageClient{})
	output, err := uc.Execute(context.Background(), updateInput(meal.ID(), otherID))

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Meal.NotFoundError))
}

func TestUpdateMeal_SaveError(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	meal := newDummyMeal(t, userID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(meal, nil)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*mealdomain.Meal")).Return(errors.New("save failed"))

	uc := mealusecase.NewUpdateMeal(repo, fakeTxManager{}, fakeStorageClient{})
	output, err := uc.Execute(context.Background(), updateInput(meal.ID(), userID))

	assert.Error(t, err)
	assert.Nil(t, output)
}
