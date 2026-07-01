package valueobject

import "errors"

type notificationTypeCode string

const (
	NotificationTypeWeeklyGoal notificationTypeCode = "weekly_goal"
)

var ErrInvalidNotificationType = errors.New("invalid notification type")

type NotificationType struct {
	LiteralBase[string]
}

func NewNotificationTypeFromString(s string) (*NotificationType, error) {
	switch notificationTypeCode(s) {
	case NotificationTypeWeeklyGoal:
		return &NotificationType{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidNotificationType
	}
}

func NewNotificationTypeFromCode(c notificationTypeCode) NotificationType {
	return NotificationType{LiteralBase: LiteralBase[string]{v: string(c)}}
}
