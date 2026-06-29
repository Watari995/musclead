package notificationpublicfunctions

import (
	"context"
	"encoding/json"

	"github.com/Watari995/musclead/internal/valueobject"
)

// NotificationCommand は goal-checker worker が通知を作成するための公開インターフェース。
type NotificationCommand interface {
	Create(ctx context.Context, userID valueobject.UserID, notificationType string, metadata json.RawMessage) error
}
