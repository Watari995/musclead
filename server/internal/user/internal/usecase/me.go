package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type MeInput struct {
	UserID valueobject.UserID
}

type MeOutput struct {
	User        userdomain.User
	Preferences userdomain.UserPreferences
}

type Me struct {
	userRepo  userdomain.UserRepository
	prefsRepo userdomain.UserPreferencesRepository
}

func (uc *Me) Execute(ctx context.Context, input MeInput) (*MeOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if user == nil {
		return nil, myerror.NewUserNotFoundError()
	}
	preferences, err := uc.prefsRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if preferences == nil {
		preferences = userdomain.CreateDefaultUserPreferences(input.UserID)
	}
	return &MeOutput{User: *user, Preferences: *preferences}, nil
}

func NewMe(userRepo userdomain.UserRepository, prefsRepo userdomain.UserPreferencesRepository) *Me {
	return &Me{userRepo: userRepo, prefsRepo: prefsRepo}
}
