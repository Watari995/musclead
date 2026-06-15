package paymentusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
)

type CreatePortalSession struct {
	paymentRepo  paymentdomain.PaymentRepository
	stripeClient paymentdomain.StripeClient
}

func (uc *CreatePortalSession) CreatePortalSession(ctx context.Context, req publicfunctions.CreatePortalSessionRequest) (publicfunctions.CreatePortalSessionResponse, error) {
	payment, err := uc.paymentRepo.FindLatestSucceededByUserID(ctx, req.UserID)
	if err != nil {
		return publicfunctions.CreatePortalSessionResponse{}, err
	}
	if payment == nil || payment.StripeCustomerID() == nil {
		return publicfunctions.CreatePortalSessionResponse{}, myerror.NewPaymentNotFoundError()
	}
	portalURL, err := uc.stripeClient.CreatePortalSession(ctx, *payment.StripeCustomerID())
	if err != nil {
		return publicfunctions.CreatePortalSessionResponse{}, err
	}
	return publicfunctions.CreatePortalSessionResponse{PortalURL: portalURL}, nil
}

func NewCreatePortalSession(paymentRepo paymentdomain.PaymentRepository, stripeClient paymentdomain.StripeClient) *CreatePortalSession {
	return &CreatePortalSession{paymentRepo: paymentRepo, stripeClient: stripeClient}
}
