package userdto

import (
	"time"

	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
)

type MeResponse struct {
	User        UserDTO        `json:"user"`
	Preferences PreferencesDTO `json:"preferences"`
}

func MeResponseFromEntities(user *userdomain.User, preferences *userdomain.UserPreferences, urlBuilder shareddomain.URLBuilder) MeResponse {
	return MeResponse{
		User:        FromEntity(user, urlBuilder),
		Preferences: PreferencesFromEntity(preferences),
	}
}

type RegisterRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Birthday *string `json:"birthday,omitempty"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type UpdateUserRequest struct {
	Name             shareddto.Patch[string] `json:"name"`
	Birthday         shareddto.Patch[string] `json:"birthday"`
	ProfileImagePath shareddto.Patch[string] `json:"profile_image_path"`
}

type UpdateUserResponse struct {
	UserID string `json:"user_id"`
}

type GenerateProfileImagePresignedURLRequest struct {
	ContentType string `json:"content_type"`
}

type GenerateProfileImagePresignedURLResponse struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type UserDTO struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Birthday        *string   `json:"birthday,omitempty"`
	ProfileImageURL string    `json:"profile_image_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func FromEntity(
	user *userdomain.User,
	urlBuilder shareddomain.URLBuilder,
) UserDTO {
	var birthdayStr *string
	if user.Birthday() != nil {
		s := user.Birthday().Format("2006-01-02")
		birthdayStr = &s
	}
	return UserDTO{
		ID:              user.ID().Value(),
		Name:            user.Name().Value(),
		Email:           user.Email().Value(),
		Birthday:        birthdayStr,
		ProfileImageURL: urlBuilder.BuildPublicURL(user.ProfileImagePath()),
		CreatedAt:       user.CreatedAt(),
		UpdatedAt:       user.UpdatedAt(),
	}
}
