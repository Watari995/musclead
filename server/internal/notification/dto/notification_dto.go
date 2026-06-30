package notificationdto

import (
	"time"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type NotificationDTO struct {
	ID               string           `json:"id"`
	NotificationType string           `json:"notification_type"`
	Metadata         valueobject.Metadata `json:"metadata"`
	IsRead           bool             `json:"is_read"`
	ReadAt           *time.Time       `json:"read_at,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
}

type GetNotificationsResponse struct {
	Notifications []NotificationDTO `json:"notifications"`
	UnreadCount   int               `json:"unread_count"`
}

func NotificationFromEntity(n *notificationdomain.Notification) NotificationDTO {
	return NotificationDTO{
		ID:               n.ID().Value(),
		NotificationType: n.NotificationType().Value(),
		Metadata:         n.Metadata(),
		IsRead:           n.IsRead(),
		ReadAt:           n.ReadAt(),
		CreatedAt:        n.CreatedAt(),
	}
}
