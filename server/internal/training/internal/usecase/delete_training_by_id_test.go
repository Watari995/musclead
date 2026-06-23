package trainingusecase_test

import (
	"context"
	"errors"
	"testing"

	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteTrainingByID_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockTrainingRepository)
	cache := new(MockExerciseBestSetTimeseriesCache)
	training := newDummyTraining(t)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(training, nil)
	repo.On("DeleteByID", mock.Anything, mock.Anything).Return(nil)

	uc := trainingusecase.NewDeleteTrainingByID(repo, cache)
	err := uc.Execute(context.Background(), trainingusecase.DeleteTrainingByIDInput{
		TrainingID: training.ID(),
		UserID:     training.UserID(),
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteTrainingByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockTrainingRepository)
	cache := new(MockExerciseBestSetTimeseriesCache)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := trainingusecase.NewDeleteTrainingByID(repo, cache)
	err := uc.Execute(context.Background(), trainingusecase.DeleteTrainingByIDInput{
		TrainingID: newDummyTraining(t).ID(),
		UserID:     newDummyTraining(t).UserID(),
	})

	assert.Error(t, err)
}

func TestDeleteTrainingByID_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockTrainingRepository)
	cache := new(MockExerciseBestSetTimeseriesCache)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	uc := trainingusecase.NewDeleteTrainingByID(repo, cache)
	err := uc.Execute(context.Background(), trainingusecase.DeleteTrainingByIDInput{
		TrainingID: newDummyTraining(t).ID(),
		UserID:     newDummyTraining(t).UserID(),
	})

	assert.Error(t, err)
}
