package main

import (
	"context"
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	notificationpublicfunctions "github.com/Watari995/musclead/internal/notification/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	"github.com/shopspring/decimal"
)

// checkAndNotify は1ユーザーの週次目標達成チェックを行い、通知を作成する。
func checkAndNotify(
	ctx context.Context,
	userID valueobject.UserID,
	weekStart time.Time,
	userQuery userpublicfunctions.UserQuery,
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
	notifCommand notificationpublicfunctions.NotificationCommand,
) error {
	// 1. 週次目標を取得
	output, err := userQuery.GetWeeklyGoal(ctx, userpublicfunctions.GetWeeklyGoalInput{UserID: userID})
	if err != nil {
		return err
	}
	// 目標が未設定なら return nil
	if output.Goal == nil {
		return nil // goalが設定されていない場合には何もしない
	}

	// 2. 各種実績を取得
	count, err := trainingQuery.CountSessionsByWeek(ctx, userID, weekStart)
	if err != nil {
		return err
	}
	calories, err := mealQuery.GetAverageCaloriesInAWeek(ctx, userID, weekStart)
	if err != nil {
		return err
	}
	change, err := weightQuery.GetWeightChangeInAWeek(ctx, userID, weekStart)
	if err != nil {
		return err
	}

	// 3. 目標と実績を比較してmetadataを組み立て
	achieved := true
	if output.Goal.TrainingCount != nil {
		achieved = achieved && count.Value() >= output.Goal.TrainingCount.Value()
	}
	if output.Goal.CalorieAverage != nil && calories != nil {
		goalCalVal := decimal.NewFromInt(int64(output.Goal.CalorieAverage.Value()))
		achieved = achieved && calories.Value().LessThanOrEqual(goalCalVal)
	}
	if output.Goal.WeightChangeKg != nil && change != nil {
		goalVal := output.Goal.WeightChangeKg.Value()
		actualVal := change.Value()
		achieved = achieved && goalVal.Sign() == actualVal.Sign() && actualVal.Abs().GreaterThanOrEqual(goalVal.Abs())
	}

	// 4. notifCommand.Create(ctx, userID, "weekly_goal", metadata)
	metadata := valueobject.Metadata{
		"training_goal":   output.Goal.TrainingCount,
		"training_actual": count,
		"calorie_goal":    output.Goal.CalorieAverage,
		"calorie_actual":  calories,
		"weight_goal":     output.Goal.WeightChangeKg,
		"weight_actual":   change,
		"achieved":        achieved,
	}

	if err := notifCommand.Create(ctx, userID, valueobject.NewNotificationTypeFromCode(valueobject.NotificationTypeWeeklyGoal), metadata); err != nil {
		return err
	}

	return nil
}
