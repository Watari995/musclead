package paymentusecase

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

// CompletePaymentInput は Webhook 'checkout.session.completed' 受信時に handler から渡される。
type CompletePaymentInput struct {
	StripeEventID string
	EventType     string
	Payload       valueobject.Metadata // event.Data の生 JSON を Metadata に詰めたもの
}

// CompletePayment は Stripe Checkout 完了時に payment を succeeded に遷移させる。
//
// 設計 (ADR 0014, 0018):
//   - TX 内で全部 atomic に実行
//   - stripe_events Create (UNIQUE 違反は ErrStripeEventAlreadyExists で no-op = 冪等性吸収)
//   - payments UPDATE (succeeded、 stripe_subscription_id, current_period_end, succeeded_at)
//   - payment_events INSERT (succeeded)
//   - outbox_events INSERT (PaymentSucceeded、 email worker 用)
type CompletePayment struct {
	paymentRepo      paymentdomain.PaymentRepository
	paymentEventRepo paymentdomain.PaymentEventRepository
	stripeEventRepo  paymentdomain.StripeEventRepository
	outboxEventRepo  paymentdomain.OutboxEventRepository
	txManager        dbtx.TransactionManager
}

// Execute は Webhook 受信時の本処理。
//
// TODO (User 実装):
//
//	return uc.txManager.Processing(ctx, func(ctx context.Context) error {
//	    // 1. stripe_events Create (冪等性吸収)
//	    stripeEvent := paymentdomain.CreateStripeEvent(input.StripeEventID, input.EventType, input.Payload)
//	    if err := uc.stripeEventRepo.Create(ctx, stripeEvent); err != nil {
//	        if errors.Is(err, paymentdomain.ErrStripeEventAlreadyExists) {
//	            return nil  // no-op、 重複受信は正常終了
//	        }
//	        return err
//	    }
//
//	    // 2. Payload から stripe_subscription_id 等を取り出す
//	    // 3. payment を引いて MarkSucceeded
//	    // 4. payment_events INSERT (succeeded)
//	    // 5. outbox INSERT (PaymentSucceeded)
//	    return nil
//	})
func (uc *CompletePayment) Execute(ctx context.Context, input CompletePaymentInput) error {
	return errNotImplemented
}

func NewCompletePayment(
	paymentRepo paymentdomain.PaymentRepository,
	paymentEventRepo paymentdomain.PaymentEventRepository,
	stripeEventRepo paymentdomain.StripeEventRepository,
	outboxEventRepo paymentdomain.OutboxEventRepository,
	txManager dbtx.TransactionManager,
) *CompletePayment {
	return &CompletePayment{
		paymentRepo:      paymentRepo,
		paymentEventRepo: paymentEventRepo,
		stripeEventRepo:  stripeEventRepo,
		outboxEventRepo:  outboxEventRepo,
		txManager:        txManager,
	}
}
