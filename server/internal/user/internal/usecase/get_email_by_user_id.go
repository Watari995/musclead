package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
)

type GetEmailByUserID struct {
	userRepo userdomain.UserRepository
}

func (uc *GetEmailByUserID) GetEmailByUserID(ctx context.Context, input publicfunctions.GetEmailByUserIDInput) (publicfunctions.GetEmailByUserIDOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return publicfunctions.GetEmailByUserIDOutput{}, myerror.NewInternalError().Wrap(err)
	}
	if user == nil {
		return publicfunctions.GetEmailByUserIDOutput{}, myerror.NewUserNotFoundError()
	}
	return publicfunctions.GetEmailByUserIDOutput{Email: user.Email()}, nil
}

func NewGetEmailByUserID(userRepo userdomain.UserRepository) *GetEmailByUserID {
	return &GetEmailByUserID{userRepo: userRepo}
}
