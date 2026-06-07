package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindUserInput struct {
	UserID valueobject.UserID
}

type FindUserOutput struct {
	User *userdomain.User
}

type FindUser struct {
	userRepo userdomain.UserRepository
}

func (uc *FindUser) Execute(ctx context.Context, input FindUserInput) (*FindUserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if user == nil {
		return nil, myerror.NewUserNotFoundError()
	}

	return &FindUserOutput{
		User: user,
	}, nil
}

func NewFindUser(userRepo userdomain.UserRepository) *FindUser {
	return &FindUser{userRepo: userRepo}
}
