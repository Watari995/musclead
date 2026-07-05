package notificationdomain

import "context"

type DeviceTokenRepository interface {
	Save(ctx context.Context, deviceToken *DeviceToken) error
}
