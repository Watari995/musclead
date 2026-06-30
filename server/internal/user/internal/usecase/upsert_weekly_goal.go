package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpsertWeeklyGoalInput struct {
	UserID         valueobject.UserID
	TrainingCount  *valueobject.NonNegativeInt
	CalorieAverage *valueobject.NonNegativeInt
	WeightChangeKg *valueobject.WeightChangeKg
}

type UpsertWeeklyGoalOutput struct {
	Goal *userdomain.UserWeeklyGoal
}

type UpsertWeeklyGoal struct {
	weeklyGoalRepo userdomain.UserWeeklyGoalRepository
}

func NewUpsertWeeklyGoal(weeklyGoalRepo userdomain.UserWeeklyGoalRepository) *UpsertWeeklyGoal {
	return &UpsertWeeklyGoal{weeklyGoalRepo: weeklyGoalRepo}
}

func (uc *UpsertWeeklyGoal) Execute(ctx context.Context, input UpsertWeeklyGoalInput) (*UpsertWeeklyGoalOutput, error) {
	goal, err := uc.weeklyGoalRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if goal == nil {
		goal = userdomain.CreateUserWeeklyGoal(
			input.UserID,
			input.TrainingCount,
			input.CalorieAverage,
			input.WeightChangeKg,
		)
	} else {
		goal.SetTrainingCount(input.TrainingCount)
		goal.SetCalorieAverage(input.CalorieAverage)
		goal.SetWeightChangeKg(input.WeightChangeKg)
	}

	if err := uc.weeklyGoalRepo.Save(ctx, goal); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpsertWeeklyGoalOutput{Goal: goal}, nil

}
