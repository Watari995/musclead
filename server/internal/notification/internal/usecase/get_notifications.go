package notificationusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
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
	notifications, err := uc.notificationRepo.FindAllByUserID(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	unreadCount := 0
	for _, notification := range notifications {
		if notification.IsRead() {
			continue
		}
		unreadCount++
	}
	return &GetNotificationsOutput{
		Notifications: notifications,
		UnreadCount:   unreadCount,
	}, nil
}
