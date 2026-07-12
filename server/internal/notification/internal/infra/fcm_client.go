package notificationinfra

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"google.golang.org/api/option"
)

type fcmClient struct {
	client *messaging.Client
}

func NewFCMClient(ctx context.Context, credentialsJSON []byte) (notificationdomain.PushNotificationClient, error) {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &fcmClient{client: client}, nil
}

func (f fcmClient) Send(ctx context.Context, token string, msg notificationdomain.PushMessage) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: msg.Title,
			Body:  msg.Body,
		},
		Data: msg.Data,
	}

	_, err := f.client.Send(ctx, message)
	if err != nil {
		if messaging.IsRegistrationTokenNotRegistered(err) {
			return notificationdomain.ErrTokenNoLongerAvailable
		} else {
			return err
		}
	}

	return nil
}
