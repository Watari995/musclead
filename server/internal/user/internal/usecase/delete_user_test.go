package userusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Watari995/musclead/internal/myerror"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUser_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)
	user := newDummyUser(t)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(nil)

	uc := userusecase.NewDeleteUser(repo)
	err := uc.Execute(context.Background(), userusecase.DeleteUserInput{UserID: user.ID()})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil)

	uc := userusecase.NewDeleteUser(repo)
	err := uc.Execute(context.Background(), userusecase.DeleteUserInput{
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.User.NotFoundError))
}

func TestDeleteUser_FindError(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("db down"))

	uc := userusecase.NewDeleteUser(repo)
	err := uc.Execute(context.Background(), userusecase.DeleteUserInput{
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
}

func TestDeleteUser_SaveError(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)
	user := newDummyUser(t)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(user, nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("save failed"))

	uc := userusecase.NewDeleteUser(repo)
	err := uc.Execute(context.Background(), userusecase.DeleteUserInput{UserID: user.ID()})

	assert.Error(t, err)
}
