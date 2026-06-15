package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
)

// userQuery は user module の Query 系 usecase を束ねて
// publicfunctions.UserQuery を満たす facade。
//
// 束ね役を別ファイル (usecase 側) に置く理由は payment の webhook_command.go のコメント参照。
type userQuery struct {
	getEmailByUserID *GetEmailByUserID
}

func NewUserQuery(getEmailByUserID *GetEmailByUserID) publicfunctions.UserQuery {
	return &userQuery{getEmailByUserID: getEmailByUserID}
}

func (q *userQuery) GetEmailByUserID(ctx context.Context, input publicfunctions.GetEmailByUserIDInput) (publicfunctions.GetEmailByUserIDOutput, error) {
	return q.getEmailByUserID.GetEmailByUserID(ctx, input)
}
