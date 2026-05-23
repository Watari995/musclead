package userusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type RegisterUserInput struct {
	Name     valueobject.String50
	Email    valueobject.Email
	Birthday *time.Time
	Password string
}

type RegisterUserOutput struct {
	UserID valueobject.UserID
}

type RegisterUser struct {
	userRepo       userdomain.UserRepository
	passwordHasher userdomain.PasswordHasher
}

func (uc *RegisterUser) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	// check email uniqueness
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil && !myerror.IsCode(err, myerror.ErrorCodes.User.NotFoundError) {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if user != nil {
		return nil, myerror.NewEmailAlreadyExistsError()
	}

	// hash password
	hash, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	// create user
	newUser := userdomain.CreateUser(input.Name, input.Email, *hash, input.Birthday)
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.User.EmailAlreadyExistsError) {
			return nil, myerror.NewEmailAlreadyExistsError()
		}
		return nil, myerror.NewInternalError().Wrap(err)
	}

	return &RegisterUserOutput{UserID: newUser.ID()}, nil
}

// NewRegisterUser is a constructor for RegisterUser
func NewRegisterUser(
	userRepo userdomain.UserRepository,
	passwordHasher userdomain.PasswordHasher,
) *RegisterUser {
	return &RegisterUser{userRepo: userRepo, passwordHasher: passwordHasher}
}
