package notificationdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type NotificationRepository interface {
	FindByIDAndUserID(ctx context.Context, id valueobject.NotificationID, userID valueobject.UserID) (*Notification, error)
	FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*Notification, error)
	Save(ctx context.Context, notification *Notification) error
	MarkAsRead(ctx context.Context, id valueobject.NotificationID) error
}
