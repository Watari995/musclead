package userusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser_Success(t *testing.T) {
	t.Parallel()

	// arrange mock
	repo := new(MockUserRepository)
	hasher := new(MockPasswordHasher)

	name, _ := valueobject.NewString50("test user")
	email, _ := valueobject.NewEmail("test@example.com")

	// bcrypt hashを作っておく
	rawHash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	hashedPassword, _ := valueobject.NewHashedPassword(string(rawHash))

	repo.On("FindByEmail", mock.Anything, *email).Return(nil, nil)
	hasher.On("Hash", "secret123").Return(hashedPassword, nil)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*userdomain.User")).Return(nil)

	// act
	uc := userusecase.NewRegisterUser(repo, hasher)
	output, err := uc.Execute(context.Background(),
		userusecase.RegisterUserInput{
			Name:     *name,
			Email:    *email,
			Password: "secret123",
		},
	)
	// assert: 期待値確認
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.UserID.Value())

	// assert: モックの検証
	repo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	t.Parallel()

	repo := new(MockUserRepository)
	hasher := new(MockPasswordHasher)

	name, _ := valueobject.NewString50("test user")
	email, _ := valueobject.NewEmail("test@example.com")

	existingUser := &userdomain.User{}

	// userを返すようにrepo.Onでmockのメソッドを設定
	repo.On("FindByEmail", mock.Anything, *email).Return(existingUser, nil)

	uc := userusecase.NewRegisterUser(repo, hasher)
	output, err := uc.Execute(context.Background(),
		userusecase.RegisterUserInput{
			Name:     *name,
			Email:    *email,
			Password: "secret123",
		},
	)
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.User.EmailAlreadyExistsError))
	repo.AssertExpectations(t)
}

func TestRegisterUser_HashPasswordError(t *testing.T) {
	t.Parallel()

	repo := new(MockUserRepository)
	hasher := new(MockPasswordHasher)

	name, _ := valueobject.NewString50("test user")
	email, _ := valueobject.NewEmail("test@example.com")

	repo.On("FindByEmail", mock.Anything, *email).Return(nil, nil)
	hasher.On("Hash", "secret123").Return(nil, errors.New("hash password error"))

	uc := userusecase.NewRegisterUser(repo, hasher)
	output, err := uc.Execute(context.Background(),
		userusecase.RegisterUserInput{
			Name:     *name,
			Email:    *email,
			Password: "secret123",
		},
	)
	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestRegisterUser_SaveError(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepository)
	hasher := new(MockPasswordHasher)

	name, _ := valueobject.NewString50("test user")
	email, _ := valueobject.NewEmail("test@example.com")
	rawHash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	hashedPassword, _ := valueobject.NewHashedPassword(string(rawHash))

	repo.On("FindByEmail", mock.Anything, *email).Return(nil, nil)
	hasher.On("Hash", "secret123").Return(hashedPassword, nil)

	repo.On("Save", mock.Anything, mock.AnythingOfType("*userdomain.User")).Return(errors.New("save error"))

	uc := userusecase.NewRegisterUser(repo, hasher)
	output, err := uc.Execute(context.Background(),
		userusecase.RegisterUserInput{
			Name:     *name,
			Email:    *email,
			Password: "secret123",
		},
	)
	assert.Error(t, err)
	assert.Nil(t, output)
}
