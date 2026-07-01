package notificationinfra

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type NotificationModel struct {
	ID               []byte               `db:"id"`
	UserID           []byte               `db:"user_id"`
	NotificationType string               `db:"notification_type"`
	Metadata         valueobject.Metadata `db:"metadata"`
	ReadAt           *time.Time           `db:"read_at"`
	CreatedAt        time.Time            `db:"created_at"`
}
