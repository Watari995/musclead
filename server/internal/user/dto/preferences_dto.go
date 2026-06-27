package userdto

import (
	"time"

	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
)

type UpdatePreferencesRequest struct {
	Theme         shareddto.Patch[string] `json:"theme"`
	MealColor     *string                 `json:"meal_color,omitempty"`
	TrainingColor *string                 `json:"training_color,omitempty"`
	WeightColor   *string                 `json:"weight_color,omitempty"`
}

type UpdatePreferencesResponse struct {
	UserID string `json:"user_id"`
}

type PreferencesDTO struct {
	Theme         string    `json:"theme"`
	MealColor     string    `json:"meal_color"`
	TrainingColor string    `json:"training_color"`
	WeightColor   string    `json:"weight_color"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func PreferencesFromEntity(p *userdomain.UserPreferences) PreferencesDTO {
	return PreferencesDTO{
		Theme:         p.Theme().Value(),
		MealColor:     p.MealColor().Value(),
		TrainingColor: p.TrainingColor().Value(),
		WeightColor:   p.WeightColor().Value(),
		CreatedAt:     p.CreatedAt(),
		UpdatedAt:     p.UpdatedAt(),
	}
}
