package main

import (
	"context"
	"log/slog"
	"sync"
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	notificationpublicfunctions "github.com/Watari995/musclead/internal/notification/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

type userCheckJob struct {
	userID    valueobject.UserID
	weekStart time.Time
}

// run はworker poolとschedulerを起動してブロックする。
func run(
	ctx context.Context,
	userQuery userpublicfunctions.UserQuery,
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
	notifCommand notificationpublicfunctions.NotificationCommand,
) {
	ch := make(chan userCheckJob, 100)

	// workerのなかで処理を行う
	var wg sync.WaitGroup
	for range 5 {
		wg.Go(func() {
			for job := range ch {
				checkAndNotify(
					ctx,
					job.userID,
					job.weekStart,
					userQuery,
					trainingQuery,
					mealQuery,
					weightQuery,
					notifCommand,
				)
			}
		})
	}

	jst := time.FixedZone("JST", 9*60*60)

	// tickerで毎週1かい処理を実行
	ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// tickerが発火した時点から7日前をweekStartとする（過去1週間を対象）。
			// JST基準で日付を計算し、UTCとの時差による日付ずれを防ぐ。
			weekStart := time.Now().In(jst).AddDate(0, 0, -7)
			go func() {
				userIDs, err := userQuery.GetAllUserIDs(ctx)
				if err != nil {
					slog.Error("failed to get user IDs", "err", err)
					return
				}
				for _, id := range userIDs {
					ch <- userCheckJob{userID: id, weekStart: weekStart}
				}
			}()
		case <-ctx.Done():
			close(ch)
			wg.Wait() // graceful shutdownのための実装
			return
		}
	}
}
