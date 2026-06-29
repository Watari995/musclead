package notificationusecase

import (
	"context"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type GetNotificationsInput struct {
	UserID valueobject.UserID
}

type GetNotificationsOutput struct {
	Notifications []*notificationdomain.Notification
	UnreadCount   int
}

type GetNotifications struct {
	notificationRepo notificationdomain.NotificationRepository
}

func NewGetNotifications(notificationRepo notificationdomain.NotificationRepository) *GetNotifications {
	return &GetNotifications{notificationRepo: notificationRepo}
}

func (uc *GetNotifications) Execute(ctx context.Context, input GetNotificationsInput) (*GetNotificationsOutput, error) {
	// TODO: implement
	return nil, nil
}
