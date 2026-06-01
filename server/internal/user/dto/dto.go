package userdto

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type UserDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Birthday  *string   `json:"birthday,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserDTO(
	id valueobject.UserID,
	name valueobject.String50,
	email valueobject.Email,
	birthday *time.Time,
	createdAt,
	updatedAt time.Time,
) UserDTO {
	var birthdayStr *string
	if birthday != nil {
		s := birthday.Format("2006-01-02")
		birthdayStr = &s
	}
	return UserDTO{
		ID:        id.Value(),
		Name:      name.Value(),
		Email:     email.Value(),
		Birthday:  birthdayStr,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
