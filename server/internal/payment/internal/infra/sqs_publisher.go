package paymentinfra

import (
	"context"
	"encoding/json"

	"github.com/Watari995/musclead/internal/myerror"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// sqsPublisher は paymentdomain.Publisher の SQS 実装。
// StripeClient と同じく、 usecase / domain は AWS SDK を直接見ない (ACL)。
// client / queueURL は Composition Root から注入する。
type sqsPublisher struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSPublisher(client *sqs.Client, queueURL string) paymentdomain.Publisher {
	return &sqsPublisher{client: client, queueURL: queueURL}
}

func (p *sqsPublisher) Publish(ctx context.Context, msg paymentdomain.PublishMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return err
	}
	return nil
}
