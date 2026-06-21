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

func TestDeleteRoutineByID_Success(t *testing.T) {
	t.Parallel()
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	repo := new(MockRoutineRepository)
	routine := newDummyRoutine(t, userID)

	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(routine, nil)
	repo.On("DeleteByID", mock.Anything, mock.Anything).Return(nil)

	uc := trainingusecase.NewDeleteRoutineByID(repo)
	err := uc.Execute(context.Background(), trainingusecase.DeleteRoutineByIDInput{
		ID:     routine.ID(),
		UserID: userID,
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteRoutineByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockRoutineRepository)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	uc := trainingusecase.NewDeleteRoutineByID(repo)
	err := uc.Execute(context.Background(), trainingusecase.DeleteRoutineByIDInput{
		ID:     valueobject.NewPrimaryID[valueobject.RoutineID](),
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
}

func TestDeleteRoutineByID_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockRoutineRepository)
	repo.On("FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	uc := trainingusecase.NewDeleteRoutineByID(repo)
	err := uc.Execute(context.Background(), trainingusecase.DeleteRoutineByIDInput{
		ID:     valueobject.NewPrimaryID[valueobject.RoutineID](),
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
}
