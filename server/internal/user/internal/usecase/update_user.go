package userusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateUserInput struct {
	UserID   valueobject.UserID
	Name     shareddto.Patch[valueobject.String50]
	Birthday shareddto.Patch[time.Time]
}

type UpdateUserOutput struct {
	UserID valueobject.UserID
}

type UpdateUser struct {
	userRepo userdomain.UserRepository
}

func (uc *UpdateUser) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		if !myerror.IsCode(err, myerror.ErrorCodes.User.NotFoundError) {
			return nil, myerror.NewInternalError().Wrap(err)
		}
		return nil, myerror.NewUserNotFoundError()
	}

	if input.Name.Set {
		user.SetName(input.Name.Value)
	}
	if input.Birthday.Set {
		var b *time.Time
		if !input.Birthday.Null {
			b = &input.Birthday.Value
		}
		user.SetBirthday(b)
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpdateUserOutput{UserID: user.ID()}, nil
}

func NewUpdateUser(userRepo userdomain.UserRepository) *UpdateUser {
	return &UpdateUser{userRepo: userRepo}
}
