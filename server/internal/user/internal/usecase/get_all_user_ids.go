package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetAllUserIDs struct {
	userRepo userdomain.UserRepository
}

func NewGetAllUserIDs(userRepo userdomain.UserRepository) *GetAllUserIDs {
	return &GetAllUserIDs{userRepo: userRepo}
}

func (uc *GetAllUserIDs) Execute(ctx context.Context) ([]valueobject.UserID, error) {
	userIDs, err := uc.userRepo.GetAllUserIDs(ctx)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return userIDs, nil
}
