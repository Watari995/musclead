package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdatePreferencesInput struct {
	UserID        valueobject.UserID
	Theme         shareddto.Patch[valueobject.Theme]
	MealColor     *valueobject.ColorHex
	TrainingColor *valueobject.ColorHex
	WeightColor   *valueobject.ColorHex
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

	// Update calendar colors
	if input.MealColor != nil {
		pref.SetMealColor(*input.MealColor)
	}
	if input.TrainingColor != nil {
		pref.SetTrainingColor(*input.TrainingColor)
	}
	if input.WeightColor != nil {
		pref.SetWeightColor(*input.WeightColor)
	}

	if err := uc.prefsRepo.Save(ctx, pref); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpdatePreferencesOutput{UserID: pref.UserID()}, nil
}

func NewUpdatePreferences(prefsRepo userdomain.UserPreferencesRepository) *UpdatePreferences {
	return &UpdatePreferences{prefsRepo: prefsRepo}
}
