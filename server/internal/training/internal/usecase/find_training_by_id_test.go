package trainingusecase_test

import (
	"context"
	"errors"
	"testing"

	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindTrainingByID_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockTrainingRepository)
	training := newDummyTraining(t)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(training, nil)

	uc := trainingusecase.NewFindTrainingByID(repo)
	out, err := uc.Execute(context.Background(), trainingusecase.FindTrainingByIDInput{
		TrainingID: training.ID(),
		UserID:     training.UserID(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, training.ID(), out.Training.ID())
	repo.AssertExpectations(t)
}

func TestFindTrainingByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockTrainingRepository)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := trainingusecase.NewFindTrainingByID(repo)
	out, err := uc.Execute(context.Background(), trainingusecase.FindTrainingByIDInput{
		TrainingID: valueobject.NewPrimaryID[valueobject.TrainingID](),
		UserID:     valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestFindTrainingByID_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockTrainingRepository)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	uc := trainingusecase.NewFindTrainingByID(repo)
	out, err := uc.Execute(context.Background(), trainingusecase.FindTrainingByIDInput{
		TrainingID: valueobject.NewPrimaryID[valueobject.TrainingID](),
		UserID:     valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestFindTrainingByID_DifferentUser_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockTrainingRepository)
	// FindByIDAndUserID returns nil when user doesn't own the training (DB-level check)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := trainingusecase.NewFindTrainingByID(repo)
	out, err := uc.Execute(context.Background(), trainingusecase.FindTrainingByIDInput{
		TrainingID: valueobject.NewPrimaryID[valueobject.TrainingID](),
		UserID:     valueobject.NewPrimaryID[valueobject.UserID](), // different user
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}
