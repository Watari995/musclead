package paymentusecase

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

// RenewPaymentInput は月次自動更新 Webhook ('invoice.payment_succeeded') の入力。
type RenewPaymentInput struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata
}

// RenewPayment は Stripe の月次自動更新成功時に payments.current_period_end を更新する。
//
// 設計 (ADR 0014, 0018):
//   - Stripe 月次課金成功時に 'invoice.payment_succeeded' Webhook 受信
//   - TX 内で stripe_events Create + payments UPDATE (current_period_end) + payment_events INSERT (renewed) + outbox INSERT (PaymentRenewed)
//   - Pro 期限 (subscriptions.expires_at) は purchase context が outbox を受けて更新
type RenewPayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

// Execute は Webhook 受信時の本処理。 CompletePayment / CancelPayment と同じパターンの TX。
//
// TODO (User 実装):
//   - 詳細は CompletePayment のコメント参照、 status は変えず current_period_end のみ更新
func (uc *RenewPayment) Execute(ctx context.Context, input RenewPaymentInput) error {
	return errNotImplemented
}

func NewRenewPayment(
	paymentRepo paymentdomain.PaymentRepository,
	paymentEventRepo paymentdomain.PaymentEventRepository,
	stripeEventRepo paymentdomain.StripeEventRepository,
	outboxEventRepo paymentdomain.OutboxEventRepository,
	txManager dbtx.TransactionManager,
) *RenewPayment {
	return &RenewPayment{
		paymentRepo:      paymentRepo,
		paymentEventRepo: paymentEventRepo,
		stripeEventRepo:  stripeEventRepo,
		outboxEventRepo:  outboxEventRepo,
		txManager:        txManager,
	}
}
