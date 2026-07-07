package notificationdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type Notification struct {
	id               valueobject.NotificationID
	userID           valueobject.UserID
	notificationType valueobject.NotificationType
	metadata         valueobject.Metadata
	readAt           *time.Time
	createdAt        time.Time
}

func (n *Notification) ID() valueobject.NotificationID {
	return n.id
}

func (n *Notification) UserID() valueobject.UserID {
	return n.userID
}

func (n *Notification) NotificationType() valueobject.NotificationType {
	return n.notificationType
}

func (n *Notification) Metadata() valueobject.Metadata {
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

func (n *Notification) ToPushMessage() (PushMessage, error) {
	var title, body string
	switch n.notificationType.Value() {
	case string(valueobject.NotificationTypeWeeklyGoal):
		title = "週間レポートが作成されました。"
		body = "アプリ右上のベルマークからご確認ください。"
	default:
		return PushMessage{}, valueobject.ErrInvalidNotificationType
	}
	return PushMessage{
		Title: title,
		Body:  body,
		Data: map[string]string{
			"notification_id": n.ID().Value(),
			"type":            n.notificationType.String(),
		},
	}, nil
}

func CreateNotification(
	userID valueobject.UserID,
	notificationType valueobject.NotificationType,
	metadata valueobject.Metadata,
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
	notificationType valueobject.NotificationType,
	metadata valueobject.Metadata,
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
