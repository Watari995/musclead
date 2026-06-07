package userusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateUserInput struct {
	UserID           valueobject.UserID
	Name             shareddto.Patch[valueobject.String50]
	Birthday         shareddto.Patch[time.Time]
	ProfileImagePath shareddto.Patch[string]
}

type UpdateUserOutput struct {
	UserID valueobject.UserID
}

type UpdateUser struct {
	userRepo      userdomain.UserRepository
	storageClient shareddomain.StorageClient
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
	var newPath string
	var oldPath string
	if input.ProfileImagePath.Set {
		if user.ProfileImagePath() != input.ProfileImagePath.Value {
			oldPath = user.ProfileImagePath()
		}
		if input.ProfileImagePath.Null {
			newPath = userdomain.DefaultProfileImagePath
		} else {
			newPath = input.ProfileImagePath.Value
		}
		user.SetProfileImagePath(newPath)
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	// if old path is not null nor default path, delete old path
	if oldPath != "" && oldPath != userdomain.DefaultProfileImagePath && oldPath != newPath {
		uc.storageClient.DeleteObject(ctx, oldPath)
	}

	return &UpdateUserOutput{UserID: user.ID()}, nil
}

func NewUpdateUser(userRepo userdomain.UserRepository, storageClient shareddomain.StorageClient) *UpdateUser {
	return &UpdateUser{userRepo: userRepo, storageClient: storageClient}
}
