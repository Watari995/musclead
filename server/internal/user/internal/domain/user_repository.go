package userdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type UserRepository interface {
	FindByID(ctx context.Context, id valueobject.UserID) (*User, error)
	FindByEmail(ctx context.Context, email valueobject.Email) (*User, error)
	Save(ctx context.Context, user *User) error
}
