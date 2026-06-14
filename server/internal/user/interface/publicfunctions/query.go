package publicfunctions

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type GetEmailByUserIDInput struct {
	UserID valueobject.UserID
}

type GetEmailByUserIDOutput struct {
	Email valueobject.Email
}

type UserQuery interface {
	GetEmailByUserID(ctx context.Context, input GetEmailByUserIDInput) (GetEmailByUserIDOutput, error)
}
