package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
)

type Authenticate struct {
	userRepo       userdomain.UserRepository
	passwordHasher userdomain.PasswordHasher
}

func (uc *Authenticate) Execute(ctx context.Context, request publicfunctions.AuthenticateRequest) (publicfunctions.AuthenticateResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return publicfunctions.AuthenticateResponse{}, myerror.NewInternalError().Wrap(err)
	}
	if user == nil {
		return publicfunctions.AuthenticateResponse{}, myerror.NewInvalidCredentialsError()
	}
	hash := user.PasswordHash()
	if err := uc.passwordHasher.Compare(request.Password, &hash); err != nil {
		return publicfunctions.AuthenticateResponse{}, myerror.NewInvalidCredentialsError()
	}

	return publicfunctions.AuthenticateResponse{UserID: user.ID()}, nil
}

func NewAuthenticate(userRepo userdomain.UserRepository, passwordHasher userdomain.PasswordHasher) *Authenticate {
	return &Authenticate{userRepo: userRepo, passwordHasher: passwordHasher}
}
