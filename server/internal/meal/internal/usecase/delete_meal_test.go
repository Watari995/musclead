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
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	meal := newDummyMeal(t, userID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(meal, nil)
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

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := mealusecase.NewDeleteMealByID(repo)
	err := uc.Execute(context.Background(), mealusecase.DeleteMealByIDInput{
		MealID: valueobject.NewPrimaryID[valueobject.MealID](),
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Meal.NotFoundError))
}

// 他ユーザーの meal を削除しようとした場合: repository が user_id で絞るので nil が返り、
// NotFound として扱われる (他人の meal の存在自体を漏らさない設計)。
func TestDeleteMealByID_OtherUserTreatedAsNotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)

	ownerID := valueobject.NewPrimaryID[valueobject.UserID]()
	otherID := valueobject.NewPrimaryID[valueobject.UserID]()
	meal := newDummyMeal(t, ownerID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := mealusecase.NewDeleteMealByID(repo)
	err := uc.Execute(context.Background(), mealusecase.DeleteMealByIDInput{
		MealID: meal.ID(),
		UserID: otherID,
	})

	assert.Error(t, err)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Meal.NotFoundError))
}

func TestDeleteMealByID_DeleteError(t *testing.T) {
	t.Parallel()
	repo := new(MockMealRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	meal := newDummyMeal(t, userID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(meal, nil)
	repo.On("DeleteByID", mock.Anything, mock.Anything).Return(errors.New("db down"))

	uc := mealusecase.NewDeleteMealByID(repo)
	err := uc.Execute(context.Background(), mealusecase.DeleteMealByIDInput{
		MealID: meal.ID(),
		UserID: userID,
	})

	assert.Error(t, err)
}
