package main

import (
	"context"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	notificationpublicfunctions "github.com/Watari995/musclead/internal/notification/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

// run はworker poolとschedulerを起動してブロックする。
// TODO: 実装する
func run(
	ctx context.Context,
	userQuery userpublicfunctions.UserQuery,
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
	notifCommand notificationpublicfunctions.NotificationCommand,
) {
	// TODO: buffered channel を作成
	// TODO: worker pool を起動
	// TODO: ticker scheduler を起動
}
