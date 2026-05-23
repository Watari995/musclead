package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type DeleteUserInput struct {
	UserID valueobject.UserID
}

type DeleteUser struct {
	userRepo userdomain.UserRepository
}

func (uc *DeleteUser) Execute(ctx context.Context, input DeleteUserInput) error {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if user == nil {
		return myerror.NewUserNotFoundError()
	}

	user.MarkAsDeleted()
	if err := uc.userRepo.Save(ctx, user); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.User.NotFoundError) {
			return myerror.NewUserNotFoundError()
		}
		return myerror.NewInternalError().Wrap(err)
	}

	return nil
}

func NewDeleteUser(userRepo userdomain.UserRepository) *DeleteUser {
	return &DeleteUser{userRepo: userRepo}
}
