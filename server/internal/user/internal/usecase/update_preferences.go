package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdatePreferencesInput struct {
	UserID valueobject.UserID
	Theme  shareddto.Patch[valueobject.Theme]
}

type UpdatePreferencesOutput struct {
	UserID valueobject.UserID
}

type UpdatePreferences struct {
	prefsRepo userdomain.UserPreferencesRepository
}

func (uc *UpdatePreferences) Execute(ctx context.Context, input UpdatePreferencesInput) (*UpdatePreferencesOutput, error) {
	pref, err := uc.prefsRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if pref == nil {
		pref = userdomain.CreateDefaultUserPreferences(input.UserID)
	}
	if input.Theme.Set {
		pref.SetTheme(input.Theme.Value)
	}
	if err := uc.prefsRepo.Save(ctx, pref); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpdatePreferencesOutput{UserID: pref.UserID()}, nil
}

func NewUpdatePreferences(prefsRepo userdomain.UserPreferencesRepository) *UpdatePreferences {
	return &UpdatePreferences{prefsRepo: prefsRepo}
}
