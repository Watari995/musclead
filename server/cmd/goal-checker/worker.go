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

// run はworker poolとschedulerを起動してブロックする。
func run(
	ctx context.Context,
	userQuery userpublicfunctions.UserQuery,
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
	notifCommand notificationpublicfunctions.NotificationCommand,
) {
	ch := make(chan valueobject.UserID, 100)

	// workerのなかで処理を行う
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userID := range ch {
				checkAndNotify(
					ctx,
					userID,
					time.Now(),
					userQuery,
					trainingQuery,
					mealQuery,
					weightQuery,
					notifCommand,
				)
			}
		}()
	}

	// tickerで毎週1かい処理を実行
	ticker := time.NewTicker(7 * 24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				userIDs, err := userQuery.GetAllUserIDs(ctx)
				if err != nil {
					slog.Error("failed to get user IDs", "err", err)
					return
				}
				for _, id := range userIDs {
					ch <- id
				}
			}()
		case <-ctx.Done():
			close(ch)
			wg.Wait() // graceful shutdownのための実装
			return
		}
	}
}
