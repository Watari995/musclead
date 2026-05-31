package publicfunctions

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type AuthenticateRequest struct {
	Email    valueobject.Email
	Password string
}

type AuthenticateResponse struct {
	UserID valueobject.UserID
}

type UserCommand interface {
	Authenticate(ctx context.Context, request AuthenticateRequest) (AuthenticateResponse, error)
}
