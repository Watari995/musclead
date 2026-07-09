package notificationdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type DeviceTokenView struct {
	ID       valueobject.DeviceTokenID
	Token    string
	Platform valueobject.NotificationPlatform
}

type DeviceTokenQuery interface {
	FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]DeviceTokenView, error)
}
