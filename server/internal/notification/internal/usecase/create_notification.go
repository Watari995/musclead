package notificationusecase

import (
	"context"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateNotificationInput struct {
	UserID           valueobject.UserID
	NotificationType valueobject.NotificationType
	Metadata         valueobject.Metadata
}

type CreateNotification struct {
	notificationRepo notificationdomain.NotificationRepository
	outboxRepo       shareddomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

func NewCreateNotification(notificationRepo notificationdomain.NotificationRepository, outboxRepo shareddomain.OutboxEventRepository, txManager dbtx.TransactionManager) *CreateNotification {
	return &CreateNotification{
		notificationRepo: notificationRepo,
		outboxRepo:       outboxRepo,
		txManager:        txManager,
	}
}

func (uc *CreateNotification) Execute(ctx context.Context, input CreateNotificationInput) error {
	// notificationを作成するときに一緒にoutbox_eventも作成する必要がある。
	n := notificationdomain.CreateNotification(input.UserID, input.NotificationType, input.Metadata)

	err := uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		err := uc.notificationRepo.Save(txCtx, n)
		if err != nil {
			return err
		}
		err = uc.outboxRepo.Save(txCtx, shareddomain.CreateOutboxEvent(
			valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypeNotification),
			n.ID().Value(),
			input.Metadata,
		))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
