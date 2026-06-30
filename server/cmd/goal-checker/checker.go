package main

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	notificationpublicfunctions "github.com/Watari995/musclead/internal/notification/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

// checkAndNotify は1ユーザーの週次目標達成チェックを行い、通知を作成する。
// TODO: 実装する
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
	//    goal, err := userQuery.GetWeeklyGoal(ctx, ...)
	//    目標が未設定なら return nil

	// 2. 各種実績を取得
	//    trainingQuery.CountSessionsByWeek(ctx, userID, weekStart)
	//    mealQuery.GetAverageCaloriesInAWeek(ctx, userID, weekStart)
	//    weightQuery.GetWeightChangeInAWeek(ctx, userID, weekStart)

	// 3. 目標と実績を比較してmetadataを組み立て

	// 4. notifCommand.Create(ctx, userID, "weekly_goal", metadata)

	return nil
}
