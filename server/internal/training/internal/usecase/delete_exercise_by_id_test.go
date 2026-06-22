package trainingusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Watari995/musclead/internal/myerror"
	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteExerciseByID_Success(t *testing.T) {
	t.Parallel()
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	repo := new(MockExerciseRepository)
	ex := newDummyExercise(t, userID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(ex, nil)
	repo.On("DeleteByID", mock.Anything, mock.Anything).Return(nil)

	uc := trainingusecase.NewDeleteExerciseByID(repo)
	err := uc.Execute(context.Background(), trainingusecase.DeleteExerciseByIDInput{
		ID:     ex.ID(),
		UserID: userID,
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteExerciseByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockExerciseRepository)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := trainingusecase.NewDeleteExerciseByID(repo)
	err := uc.Execute(context.Background(), trainingusecase.DeleteExerciseByIDInput{
		ID:     valueobject.NewPrimaryID[valueobject.ExerciseID](),
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
}

func TestDeleteExerciseByID_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockExerciseRepository)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	uc := trainingusecase.NewDeleteExerciseByID(repo)
	err := uc.Execute(context.Background(), trainingusecase.DeleteExerciseByIDInput{
		ID:     valueobject.NewPrimaryID[valueobject.ExerciseID](),
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
}

func TestDeleteExerciseByID_UsedInTraining(t *testing.T) {
	t.Parallel()
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	repo := new(MockExerciseRepository)
	ex := newDummyExercise(t, userID)
	usedErr := myerror.NewExerciseUsedInTrainingError()

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(ex, nil)
	repo.On("DeleteByID", mock.Anything, mock.Anything).Return(usedErr)

	uc := trainingusecase.NewDeleteExerciseByID(repo)
	err := uc.Execute(context.Background(), trainingusecase.DeleteExerciseByIDInput{
		ID:     ex.ID(),
		UserID: userID,
	})

	assert.Error(t, err)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Training.ExerciseUsedInTrainingError))
	repo.AssertExpectations(t)
}
