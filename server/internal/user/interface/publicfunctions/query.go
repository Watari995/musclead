package publicfunctions

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type GetEmailByUserIDInput struct {
	UserID valueobject.UserID
}

type GetEmailByUserIDOutput struct {
	Email valueobject.Email
}

type GetWeeklyGoalInput struct {
	UserID valueobject.UserID
}

type GetWeeklyGoalOutput struct {
	ID             valueobject.UserWeeklyGoalID
	UserID         valueobject.UserID
	TrainingCount  *valueobject.NonNegativeInt
	CalorieAverage *valueobject.NonNegativeInt
	WeightChangeKg *valueobject.WeightChangeKg
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserQuery interface {
	GetEmailByUserID(ctx context.Context, input GetEmailByUserIDInput) (GetEmailByUserIDOutput, error)
	GetWeeklyGoal(ctx context.Context, input GetWeeklyGoalInput) (GetWeeklyGoalOutput, error)
	GetAllUserIDs(ctx context.Context) ([]valueobject.UserID, error)
}
