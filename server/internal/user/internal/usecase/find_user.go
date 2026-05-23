package userusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type FindUserInput struct {
	UserID valueobject.UserID
}

type FindUserOutput struct {
	UserID    valueobject.UserID
	Name      valueobject.String50
	Email     valueobject.Email
	Birthday  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
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
		UserID:    user.ID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Birthday:  user.Birthday(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}

func NewFindUser(userRepo userdomain.UserRepository) *FindUser {
	return &FindUser{userRepo: userRepo}
}
