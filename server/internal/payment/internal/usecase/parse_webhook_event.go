package paymentusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
)

type ParseWebhookEvent struct {
	stripeClient paymentdomain.StripeClient
}

func (uc *ParseWebhookEvent) ParseAndVerify(ctx context.Context, payload []byte, signatureHeader string) (publicfunctions.StripeEvent, error) {
	output, err := uc.stripeClient.ParseWebhookEvent(ctx, paymentdomain.ParseWebhookEventInput{
		Payload:         payload,
		SignatureHeader: signatureHeader,
	})
	if err != nil {
		return publicfunctions.StripeEvent{}, err
	}
	return publicfunctions.StripeEvent{
		StripeEventID: output.StripeEventID,
		EventType:     output.EventType,
		Payload:       output.Payload,
	}, nil
}

func NewParseWebhookEvent(stripeClient paymentdomain.StripeClient) *ParseWebhookEvent {
	return &ParseWebhookEvent{stripeClient: stripeClient}
}
