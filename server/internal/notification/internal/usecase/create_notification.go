package notificationusecase

import (
	"context"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateNotificationInput struct {
	UserID           valueobject.UserID
	NotificationType valueobject.NotificationType
	Metadata         valueobject.Metadata
}

type CreateNotification struct {
	notificationRepo notificationdomain.NotificationRepository
}

func NewCreateNotification(notificationRepo notificationdomain.NotificationRepository) *CreateNotification {
	return &CreateNotification{notificationRepo: notificationRepo}
}

func (uc *CreateNotification) Execute(ctx context.Context, input CreateNotificationInput) error {
	if err := uc.notificationRepo.Save(ctx, notificationdomain.CreateNotification(
		input.UserID,
		input.NotificationType,
		input.Metadata,
	)); err != nil {
		return err
	}
	return nil
}
