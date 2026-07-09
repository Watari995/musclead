package notificationdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type DeviceTokenRepository interface {
	Save(ctx context.Context, deviceToken *DeviceToken) error
	DeleteByID(ctx context.Context, id valueobject.DeviceTokenID) error
}
