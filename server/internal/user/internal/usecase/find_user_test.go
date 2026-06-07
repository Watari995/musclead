package userusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func newDummyUser(t *testing.T) *userdomain.User {
	t.Helper()
	id := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("dummy")
	email, _ := valueobject.NewEmail("dummy@example.com")
	rawHash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	hashed, _ := valueobject.NewHashedPassword(string(rawHash))
	profileImagePath := "dummy/profile.png"
	now := time.Now()
	return userdomain.NewUser(id, *name, *email, *hashed, nil, profileImagePath, now, now)
}

func TestFindUser_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)
	user := newDummyUser(t)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(user, nil)

	uc := userusecase.NewFindUser(repo)
	output, err := uc.Execute(context.Background(), userusecase.FindUserInput{UserID: user.ID()})

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, user.ID().Value(), output.User.ID().Value())
	repo.AssertExpectations(t)
}

func TestFindUser_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil)

	uc := userusecase.NewFindUser(repo)
	output, err := uc.Execute(context.Background(), userusecase.FindUserInput{
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.User.NotFoundError))
}

func TestFindUser_DBError(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)

	repo.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("db down"))

	uc := userusecase.NewFindUser(repo)
	output, err := uc.Execute(context.Background(), userusecase.FindUserInput{
		UserID: valueobject.NewPrimaryID[valueobject.UserID](),
	})

	assert.Error(t, err)
	assert.Nil(t, output)
}
