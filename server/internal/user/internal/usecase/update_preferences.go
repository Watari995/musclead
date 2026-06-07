package userusecase

import (
	"context"

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
	return nil, nil
}

func NewUpdatePreferences(prefsRepo userdomain.UserPreferencesRepository) *UpdatePreferences {
	return &UpdatePreferences{prefsRepo: prefsRepo}
}
