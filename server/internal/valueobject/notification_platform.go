package valueobject

import "errors"

type notificationPlatformCode string

const (
	NotificationPlatformIOS     notificationPlatformCode = "ios"
	NotificationPlatformAndroid notificationPlatformCode = "android"
)

var ErrInvalidNotificationPlatform = errors.New("invalid notification platform")

type NotificationPlatform struct {
	LiteralBase[string]
}

func NewNotificationPlatformFromString(s string) (*NotificationPlatform, error) {
	switch s {
	case string(NotificationPlatformIOS), string(NotificationPlatformAndroid):
		return &NotificationPlatform{
			LiteralBase[string]{v: s},
		}, nil
	default:
		return nil, ErrInvalidNotificationPlatform
	}
}

func NewNotificationPlatformFromCode(c notificationPlatformCode) NotificationPlatform {
	return NotificationPlatform{
		LiteralBase[string]{v: string(c)},
	}
}
