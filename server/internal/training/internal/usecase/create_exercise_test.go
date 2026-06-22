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

func TestCreateExercise_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockExerciseRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Squat")

	repo.On("NextDisplayOrder", mock.Anything, mock.Anything).Return(0, nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(nil)

	uc := trainingusecase.NewCreateExercise(repo)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateExerciseInput{
		UserID: userID,
		Name:   *name,
	})

	assert.NoError(t, err)
	assert.NotNil(t, out)
	repo.AssertExpectations(t)
}

func TestCreateExercise_NameAlreadyExists(t *testing.T) {
	t.Parallel()
	repo := new(MockExerciseRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Squat")

	repo.On("NextDisplayOrder", mock.Anything, mock.Anything).Return(0, nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(myerror.NewExerciseNameAlreadyExistsError())

	uc := trainingusecase.NewCreateExercise(repo)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateExerciseInput{
		UserID: userID,
		Name:   *name,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Training.ExerciseNameAlreadyExistsError))
}

func TestCreateExercise_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockExerciseRepository)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Deadlift")

	repo.On("NextDisplayOrder", mock.Anything, mock.Anything).Return(0, errors.New("db error"))

	uc := trainingusecase.NewCreateExercise(repo)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateExerciseInput{
		UserID: userID,
		Name:   *name,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}
