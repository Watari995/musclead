package purchaseusecase

import (
	"context"

	paymentpublicfunctions "github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreatePortalSession struct {
	paymentCommand paymentpublicfunctions.PaymentCommand
}

type CreatePortalSessionInput struct {
	UserID valueobject.UserID
}

type CreatePortalSessionOutput struct {
	PortalURL valueobject.URL
}

func (uc *CreatePortalSession) Execute(ctx context.Context, input CreatePortalSessionInput) (CreatePortalSessionOutput, error) {
	output, err := uc.paymentCommand.CreatePortalSession(ctx, paymentpublicfunctions.CreatePortalSessionRequest{UserID: input.UserID})
	if err != nil {
		return CreatePortalSessionOutput{}, err
	}
	return CreatePortalSessionOutput{PortalURL: output.PortalURL}, nil
}

func NewCreatePortalSession(paymentCommand paymentpublicfunctions.PaymentCommand) *CreatePortalSession {
	return &CreatePortalSession{paymentCommand: paymentCommand}
}
