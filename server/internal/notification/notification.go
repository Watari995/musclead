// Package notification は notification モジュールの公開 Facade。
package notification

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	notificationpublicfunctions "github.com/Watari995/musclead/internal/notification/interface/publicfunctions"
	notificationhandler "github.com/Watari995/musclead/internal/notification/internal/handler"
	notificationinfra "github.com/Watari995/musclead/internal/notification/internal/infra"
	notificationusecase "github.com/Watari995/musclead/internal/notification/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	outboxinfra "github.com/Watari995/musclead/internal/shared/infra/outbox"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	Handler             http.Handler
	notificationCommand notificationpublicfunctions.NotificationCommand
	relayOutbox         *notificationusecase.RelayOutbox
	relayEnabled        bool // CredentialsJSONが設定されている時だけrelayを回す。
}

type Config struct {
	CredentialsJSON []byte
}

func NewModule(dbmap *gorp.DbMap, cfg Config) *Module {
	dbmap.AddTableWithName(notificationinfra.NotificationModel{}, "notifications").SetKeys(false, "ID")
	dbmap.AddTableWithName(notificationinfra.DeviceTokenModel{}, "device_tokens").SetKeys(false, "ID")
	repo := notificationinfra.NewNotificationRepository(dbmap)
	deviceTokenRepo := notificationinfra.NewDeviceTokenRepository(dbmap)
	deviceTokenQuery := notificationinfra.NewDeviceTokenQuery(dbmap)
	outboxEventRepo := outboxinfra.NewOutboxEventRepository(dbmap)
	fcmClient, err := notificationinfra.NewFCMClient(context.Background(), cfg.CredentialsJSON)
	if err != nil {
		slog.Error("Initialize push notification client", "err", err)
	}
	txManager := dbtx.NewTransactionManager(dbmap)

	getNotifications := notificationusecase.NewGetNotifications(repo)
	getNotification := notificationusecase.NewGetNotification(repo)
	readNotification := notificationusecase.NewReadNotification(repo)
	registerDeviceToken := notificationusecase.NewRegisterDeviceToken(deviceTokenRepo)

	createNotification := notificationusecase.NewCreateNotification(repo, outboxEventRepo, txManager)

	var relayOutbox *notificationusecase.RelayOutbox
	relayEnabled := true
	if len(cfg.CredentialsJSON) != 0 && err == nil {
		relayOutbox = notificationusecase.NewRelayOutbox(outboxEventRepo, repo, deviceTokenQuery, deviceTokenRepo, fcmClient)
	} else {
		relayEnabled = false
	}

	return &Module{
		Handler:             notificationhandler.New(getNotifications, getNotification, readNotification, registerDeviceToken),
		notificationCommand: &notificationCommand{create: createNotification},
		relayOutbox:         relayOutbox,
		relayEnabled:        relayEnabled,
	}
}

func (m *Module) NotificationCommand() notificationpublicfunctions.NotificationCommand {
	return m.notificationCommand
}

func (m *Module) RunRelay(ctx context.Context) {
	if !m.relayEnabled {
		slog.Info("push notification relay disabled (no FCM credentials)")
		return
	}

	// 10秒ごとにポーリングをしてpendingのnotificationを使用してpush通知を送る
	relayInterval := 10 * time.Second
	ticker := time.NewTicker(relayInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.relayOutbox.Execute(ctx); err != nil {
				slog.Error("outbox relay failed", "err", err)
			}
		}
	}
}

// notificationCommand は NotificationCommand インターフェースの内部実装。
type notificationCommand struct {
	create *notificationusecase.CreateNotification
}

func (c *notificationCommand) Create(ctx context.Context, userID valueobject.UserID, notificationType valueobject.NotificationType, metadata valueobject.Metadata) error {
	return c.create.Execute(ctx, notificationusecase.CreateNotificationInput{
		UserID:           userID,
		NotificationType: notificationType,
		Metadata:         metadata,
	})
}
