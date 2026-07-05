package notificationusecase

import (
	"context"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type RegisterDeviceTokenInput struct {
	UserID   valueobject.UserID
	Token    string
	Platform valueobject.NotificationPlatform
}

type RegisterDeviceToken struct {
	deviceTokenRepo notificationdomain.DeviceTokenRepository
}

func NewRegisterDeviceToken(deviceTokenRepo notificationdomain.DeviceTokenRepository) *RegisterDeviceToken {
	return &RegisterDeviceToken{deviceTokenRepo: deviceTokenRepo}
}

func (uc *RegisterDeviceToken) Execute(ctx context.Context, input RegisterDeviceTokenInput) error {
	return uc.deviceTokenRepo.Save(ctx, notificationdomain.CreateDeviceToken(
		input.UserID,
		input.Token,
		input.Platform,
	))
}
