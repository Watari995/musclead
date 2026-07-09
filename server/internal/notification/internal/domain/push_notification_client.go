package notificationdomain

import (
	"context"
	"errors"
)

type PushMessage struct {
	Title string
	Body  string
	Data  map[string]string // タップ時の遷移先など、アプリ側で使うカスタムデータ
}

var ErrTokenNoLongerAvailable = errors.New("push notification token is no longer available")

type PushNotificationClient interface {
	Send(ctx context.Context, token string, msg PushMessage) error
}
