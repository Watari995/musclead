package userusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

// userQuery は user module の Query 系 usecase を束ねて
// publicfunctions.UserQuery を満たす facade。
//
// 束ね役を別ファイル (usecase 側) に置く理由は payment の webhook_command.go のコメント参照。
type userQuery struct {
	getEmailByUserID *GetEmailByUserID
	getAllUserIDs    *GetAllUserIDs
	getWeeklyGoal    *GetWeeklyGoal
}

func NewUserQuery(
	getEmailByUserID *GetEmailByUserID,
	getAllUserIDs *GetAllUserIDs,
	getWeeklyGoal *GetWeeklyGoal,
) publicfunctions.UserQuery {
	return &userQuery{
		getEmailByUserID: getEmailByUserID,
		getAllUserIDs:    getAllUserIDs,
		getWeeklyGoal:    getWeeklyGoal,
	}
}

func (q *userQuery) GetEmailByUserID(ctx context.Context, input publicfunctions.GetEmailByUserIDInput) (publicfunctions.GetEmailByUserIDOutput, error) {
	return q.getEmailByUserID.GetEmailByUserID(ctx, input)
}

func (q *userQuery) GetWeeklyGoal(ctx context.Context, input publicfunctions.GetWeeklyGoalInput) (publicfunctions.GetWeeklyGoalOutput, error) {
	output, err := q.getWeeklyGoal.Execute(ctx, GetWeeklyGoalInput{UserID: input.UserID})
	if err != nil {
		return publicfunctions.GetWeeklyGoalOutput{}, err
	}
	return publicfunctions.GetWeeklyGoalOutput{
		ID:             output.Goal.ID(),
		UserID:         output.Goal.UserID(),
		TrainingCount:  output.Goal.TrainingCount(),
		CalorieAverage: output.Goal.CalorieAverage(),
		WeightChangeKg: output.Goal.WeightChangeKg(),
		CreatedAt:      output.Goal.CreatedAt(),
		UpdatedAt:      output.Goal.UpdatedAt(),
	}, nil
}

func (q *userQuery) GetAllUserIDs(ctx context.Context) ([]valueobject.UserID, error) {
	return q.getAllUserIDs.Execute(ctx)
}
