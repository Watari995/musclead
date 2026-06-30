package main

import (
	"context"
	"sync"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	notificationpublicfunctions "github.com/Watari995/musclead/internal/notification/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

const (
	workerCount = 5
	// 本番: 7 * 24 * time.Hour
	// テスト時は短くして動作確認できる
	checkInterval = 7 * 24 * time.Hour
)

// run はworker poolとschedulerを起動してブロックする。
// ctxがキャンセルされると全workerを停止して返る。
func run(
	ctx context.Context,
	userQuery userpublicfunctions.UserQuery,
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
	notifCommand notificationpublicfunctions.NotificationCommand,
) {
	ch := make(chan valueobject.UserID, workerCount*2)

	// worker pool
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userID := range ch {
				weekStart := lastSunday()
				if err := checkAndNotify(ctx, userID, weekStart, userQuery, trainingQuery, mealQuery, weightQuery, notifCommand); err != nil {
					// TODO: structured logging
					_ = err
				}
			}
		}()
	}

	// scheduler
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			userIDs, err := userQuery.GetAllUserIDs(ctx)
			if err != nil {
				_ = err
				continue
			}
			for _, id := range userIDs {
				ch <- id
			}
		case <-ctx.Done():
			close(ch)
			wg.Wait()
			return
		}
	}
}

// lastSunday は直近の日曜日（週の起点）を JST で返す。
func lastSunday() time.Time {
	jst := time.FixedZone("JST", 9*60*60)
	now := time.Now().In(jst)
	daysBack := int(now.Weekday())
	return time.Date(now.Year(), now.Month(), now.Day()-daysBack, 0, 0, 0, 0, jst)
}
