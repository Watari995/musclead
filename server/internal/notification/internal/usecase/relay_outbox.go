package notificationusecase

import (
	"context"
	"errors"
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
	deviceTokenRepo        notificationdomain.DeviceTokenRepository
	pushNotificationClient notificationdomain.PushNotificationClient
}

func NewRelayOutbox(outboxRepo shareddomain.OutboxEventRepository, notificationRepo notificationdomain.NotificationRepository, deviceTokenQuery notificationdomain.DeviceTokenQuery, deviceTokenRepo notificationdomain.DeviceTokenRepository, pushNotificationClient notificationdomain.PushNotificationClient) *RelayOutbox {
	return &RelayOutbox{
		outboxRepo:             outboxRepo,
		notificationRepo:       notificationRepo,
		deviceTokenQuery:       deviceTokenQuery,
		deviceTokenRepo:        deviceTokenRepo,
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
		return nil // 処理するイベントがなかったら何もしない。
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
				if errors.Is(err, notificationdomain.ErrTokenNoLongerAvailable) {
					// tokenが無効だった場合はdevice_tokenから消す。
					if err := uc.deviceTokenRepo.DeleteByID(ctx, t.ID); err != nil {
						return err
					}
					slog.Info("relay: removed stale device token", "device_token_id", t.ID.Value())
					continue
				}
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
