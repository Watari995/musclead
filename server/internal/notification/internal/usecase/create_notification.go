package notificationusecase

import (
	"context"
	"encoding/json"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateNotificationInput struct {
	UserID           valueobject.UserID
	NotificationType string
	Metadata         json.RawMessage
}

type CreateNotification struct {
	notificationRepo notificationdomain.NotificationRepository
}

func NewCreateNotification(notificationRepo notificationdomain.NotificationRepository) *CreateNotification {
	return &CreateNotification{notificationRepo: notificationRepo}
}

func (uc *CreateNotification) Execute(ctx context.Context, input CreateNotificationInput) error {
	// TODO: implement
	return nil
}
