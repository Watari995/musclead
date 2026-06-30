package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetWeeklyGoalInput struct {
	UserID valueobject.UserID
}

type GetWeeklyGoalOutput struct {
	Goal *userdomain.UserWeeklyGoal
}

type GetWeeklyGoal struct {
	weeklyGoalRepo userdomain.UserWeeklyGoalRepository
}

func NewGetWeeklyGoal(weeklyGoalRepo userdomain.UserWeeklyGoalRepository) *GetWeeklyGoal {
	return &GetWeeklyGoal{weeklyGoalRepo: weeklyGoalRepo}
}

func (uc *GetWeeklyGoal) Execute(ctx context.Context, input GetWeeklyGoalInput) (*GetWeeklyGoalOutput, error) {
	goal, err := uc.weeklyGoalRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	return &GetWeeklyGoalOutput{Goal: goal}, nil
}
