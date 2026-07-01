package userdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type UserWeeklyGoalRepository interface {
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*UserWeeklyGoal, error)
	Save(ctx context.Context, goal *UserWeeklyGoal) error
}
