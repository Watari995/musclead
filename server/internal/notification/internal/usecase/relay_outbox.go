package notificationusecase

import (
	"context"
	"log/slog"

	"github.com/Watari995/musclead/internal/myerror"
	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// relayBatchSizeは1回のポーリングで処理するoutboxの最大件数。
const relayBatchSize = 50

type RelayOutbox struct {
	outboxRepo             shareddomain.OutboxEventRepository
	notificationRepo       notificationdomain.NotificationRepository
	deviceTokenQuery       notificationdomain.DeviceTokenQuery
	pushNotificationClient notificationdomain.PushNotificationClient
}

func NewRelayOutBox(outboxRepo shareddomain.OutboxEventRepository, notificationRepo notificationdomain.NotificationRepository, deviceTokenQuery notificationdomain.DeviceTokenQuery, pushNotificationClient notificationdomain.PushNotificationClient) *RelayOutbox {
	return &RelayOutbox{
		outboxRepo:             outboxRepo,
		notificationRepo:       notificationRepo,
		deviceTokenQuery:       deviceTokenQuery,
		pushNotificationClient: pushNotificationClient,
	}
}

func (uc *RelayOutbox) Execute(ctx context.Context) error {
	events, err := uc.outboxRepo.FindPendingByEventTypes(ctx, []valueobject.OutboxEventType{
		valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypeNotification),
	}, relayBatchSize)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if len(events) == 0 {
		return nil // 処理するエラーがなかったら何もしない。
	}

	// notificationを1件ずつ処理していく
	for _, e := range events {
		notificationID, err := valueobject.NewPrimaryIDFromString[valueobject.NotificationID](e.AggregateID())
		if err != nil {
			slog.Error("relay: failed to publish", "err", err)
			continue
		}
		notification, err := uc.notificationRepo.FindByID(ctx, *notificationID)
		if err != nil {
			slog.Error("relay: failed to publish", "err", err)
			continue
		}
		if notification == nil {
			slog.Error("relay: failed to publish because notification is nil")
			continue
		}
		// userIDからtokenを取得
		tokens, err := uc.deviceTokenQuery.FindAllByUserID(ctx, notification.UserID())
		if err != nil {
			slog.Error("relay: failed to publish", "err", err)
			continue
		}
		msg, err := notification.ToPushMessage()
		if err != nil {
			slog.Error("relay: failed to publish", "err", err)
			continue
		}

		for _, t := range tokens {
			// 送信処理(best effort)
			if err := uc.pushNotificationClient.Send(ctx, t.Token, msg); err != nil {
				slog.Error("relay: failed to publish", "err", err)
				continue
			}
		}

		// 成功とマークする
		e.MarkPublished()
		if err := uc.outboxRepo.Save(ctx, e); err != nil {
			return myerror.NewInternalError().Wrap(err)
		}
	}
	return nil
}
