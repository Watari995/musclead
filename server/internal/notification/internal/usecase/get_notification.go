package notificationusecase

import (
	"context"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetNotificationInput struct {
	ID     valueobject.NotificationID
	UserID valueobject.UserID
}

type GetNotificationOutput struct {
	Notification *notificationdomain.Notification
}

type GetNotification struct {
	notificationRepo notificationdomain.NotificationRepository
}

func NewGetNotification(notificationRepo notificationdomain.NotificationRepository) *GetNotification {
	return &GetNotification{notificationRepo: notificationRepo}
}

func (uc *GetNotification) Execute(ctx context.Context, input GetNotificationInput) (*GetNotificationOutput, error) {
	// TODO: implement
	return nil, nil
}
