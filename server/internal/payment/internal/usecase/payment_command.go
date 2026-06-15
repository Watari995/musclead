package paymentusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/payment/interface/publicfunctions"
)

// paymentCommand は payment module の Command 系 usecase を束ねて
// publicfunctions.PaymentCommand を満たす facade。
//
// 束ね役を module ファイル (payment.go) ではなく usecase 側の別ファイルに置く理由は
// webhook_command.go のコメント参照 (module は薄い Composition Root に保つ / 委譲ロジックは usecase の隣)。
type paymentCommand struct {
	initiatePayment     *InitiatePayment
	createPortalSession *CreatePortalSession
}

func NewPaymentCommand(initiatePayment *InitiatePayment, createPortalSession *CreatePortalSession) publicfunctions.PaymentCommand {
	return &paymentCommand{initiatePayment: initiatePayment, createPortalSession: createPortalSession}
}

func (c *paymentCommand) InitiatePayment(ctx context.Context, req publicfunctions.InitiatePaymentRequest) (publicfunctions.InitiatePaymentResponse, error) {
	return c.initiatePayment.InitiatePayment(ctx, req)
}

func (c *paymentCommand) CreatePortalSession(ctx context.Context, req publicfunctions.CreatePortalSessionRequest) (publicfunctions.CreatePortalSessionResponse, error) {
	return c.createPortalSession.CreatePortalSession(ctx, req)
}
