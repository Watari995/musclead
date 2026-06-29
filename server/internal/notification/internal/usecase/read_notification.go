package notificationusecase

import (
	"context"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ReadNotificationInput struct {
	ID     valueobject.NotificationID
	UserID valueobject.UserID
}

type ReadNotificationOutput struct{}

type ReadNotification struct {
	notificationRepo notificationdomain.NotificationRepository
}

func NewReadNotification(notificationRepo notificationdomain.NotificationRepository) *ReadNotification {
	return &ReadNotification{notificationRepo: notificationRepo}
}

func (uc *ReadNotification) Execute(ctx context.Context, input ReadNotificationInput) (*ReadNotificationOutput, error) {
	// TODO: implement
	return nil, nil
}
