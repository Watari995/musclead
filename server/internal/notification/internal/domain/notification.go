package notificationdomain

import (
	"encoding/json"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type Notification struct {
	id               valueobject.NotificationID
	userID           valueobject.UserID
	notificationType string
	metadata         json.RawMessage
	readAt           *time.Time
	createdAt        time.Time
}

func (n *Notification) ID() valueobject.NotificationID {
	return n.id
}

func (n *Notification) UserID() valueobject.UserID {
	return n.userID
}

func (n *Notification) NotificationType() string {
	return n.notificationType
}

func (n *Notification) Metadata() json.RawMessage {
	return n.metadata
}

func (n *Notification) ReadAt() *time.Time {
	return n.readAt
}

func (n *Notification) IsRead() bool {
	return n.readAt != nil
}

func (n *Notification) CreatedAt() time.Time {
	return n.createdAt
}

func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.readAt = &now
}

func CreateNotification(
	userID valueobject.UserID,
	notificationType string,
	metadata json.RawMessage,
) *Notification {
	return &Notification{
		id:               valueobject.NewPrimaryID[valueobject.NotificationID](),
		userID:           userID,
		notificationType: notificationType,
		metadata:         metadata,
		createdAt:        time.Now(),
	}
}

func NewNotification(
	id valueobject.NotificationID,
	userID valueobject.UserID,
	notificationType string,
	metadata json.RawMessage,
	readAt *time.Time,
	createdAt time.Time,
) *Notification {
	return &Notification{
		id:               id,
		userID:           userID,
		notificationType: notificationType,
		metadata:         metadata,
		readAt:           readAt,
		createdAt:        createdAt,
	}
}
