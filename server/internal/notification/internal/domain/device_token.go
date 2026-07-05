package notificationdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type DeviceToken struct {
	id        valueobject.DeviceTokenID
	userID    valueobject.UserID
	token     string
	platform  valueobject.NotificationPlatform
	createdAt time.Time
	updatedAt time.Time
}

func (d *DeviceToken) ID() valueobject.DeviceTokenID {
	return d.id
}

func (d *DeviceToken) UserID() valueobject.UserID {
	return d.userID
}

func (d *DeviceToken) Token() string {
	return d.token
}

func (d *DeviceToken) Platform() valueobject.NotificationPlatform {
	return d.platform
}

func (d *DeviceToken) CreatedAt() time.Time {
	return d.createdAt
}

func (d *DeviceToken) UpdatedAt() time.Time {
	return d.updatedAt
}

func CreateDeviceToken(
	id valueobject.DeviceTokenID,
	userID valueobject.UserID,
	token string,
	platform valueobject.NotificationPlatform,
) *DeviceToken {
	now := time.Now()
	return &DeviceToken{
		id:        valueobject.NewPrimaryID[valueobject.DeviceTokenID](),
		userID:    userID,
		token:     token,
		platform:  platform,
		createdAt: now,
		updatedAt: now,
	}
}

func NewDeviceToken(
	id valueobject.DeviceTokenID,
	userID valueobject.UserID,
	token string,
	platform valueobject.NotificationPlatform,
	createdAt time.Time,
	updatedAt time.Time,
) *DeviceToken {
	return &DeviceToken{
		id:        id,
		userID:    userID,
		token:     token,
		platform:  platform,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
