package userdto

import (
	"time"

	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
)

type UpdatePreferencesRequest struct {
	Theme shareddto.Patch[string] `json:"theme"`
}

type UpdatePreferencesResponse struct {
	UserID string `json:"user_id"`
}

type PreferencesDTO struct {
	Theme     string    `json:"theme"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func PreferencesFromEntity(p *userdomain.UserPreferences) PreferencesDTO {
	return PreferencesDTO{
		Theme:     p.Theme().Value(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}
