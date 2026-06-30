package notificationusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
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
	notification, err := uc.notificationRepo.FindByIDAndUserID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if notification == nil {
		return nil, myerror.NewNotificationNotFoundError()
	}

	return &GetNotificationOutput{
		Notification: notification,
	}, nil
}
