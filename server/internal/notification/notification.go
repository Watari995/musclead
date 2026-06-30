// Package notification は notification モジュールの公開 Facade。
package notification

import (
	"context"
	"net/http"

	notificationpublicfunctions "github.com/Watari995/musclead/internal/notification/interface/publicfunctions"
	notificationhandler "github.com/Watari995/musclead/internal/notification/internal/handler"
	notificationinfra "github.com/Watari995/musclead/internal/notification/internal/infra"
	notificationusecase "github.com/Watari995/musclead/internal/notification/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	Handler             http.Handler
	notificationCommand notificationpublicfunctions.NotificationCommand
}

func NewModule(dbmap *gorp.DbMap) *Module {
	dbmap.AddTableWithName(notificationinfra.NotificationModel{}, "notifications").SetKeys(false, "ID")
	repo := notificationinfra.NewNotificationRepository(dbmap)

	getNotifications := notificationusecase.NewGetNotifications(repo)
	getNotification := notificationusecase.NewGetNotification(repo)
	readNotification := notificationusecase.NewReadNotification(repo)

	createNotification := notificationusecase.NewCreateNotification(repo)

	return &Module{
		Handler:             notificationhandler.New(getNotifications, getNotification, readNotification),
		notificationCommand: &notificationCommand{create: createNotification},
	}
}

func (m *Module) NotificationCommand() notificationpublicfunctions.NotificationCommand {
	return m.notificationCommand
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
