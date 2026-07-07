package notificationdomain

import "context"

type PushMessage struct {
	Title string
	Body  string
	Data  map[string]string // タップ時の遷移先など、アプリ側で使うカスタムデータ
}

type PushNotificationClient interface {
	Send(ctx context.Context, token string, msg PushMessage) error
}
