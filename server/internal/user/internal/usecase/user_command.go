package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
)

// userCommand は user module の Command 系 usecase を束ねて
// publicfunctions.UserCommand を満たす facade。
//
// 束ね役を別ファイル (usecase 側) に置く理由は payment の webhook_command.go のコメント参照。
type userCommand struct {
	authenticate *Authenticate
}

func NewUserCommand(authenticate *Authenticate) publicfunctions.UserCommand {
	return &userCommand{authenticate: authenticate}
}

func (c *userCommand) Authenticate(ctx context.Context, request publicfunctions.AuthenticateRequest) (publicfunctions.AuthenticateResponse, error) {
	return c.authenticate.Authenticate(ctx, request)
}
