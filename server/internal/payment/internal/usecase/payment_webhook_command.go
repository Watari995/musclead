package paymentusecase

// PaymentWebhookCommand は 4 つの Webhook 系 usecase を 1 つの struct に集約し、
// publicfunctions.PaymentWebhookCommand interface を満たすためのラッパー。
//
// 設計 (ADR 0019):
//   - billing.NewModule に渡せる 1 つの依存に集約する役割
//   - struct embedding で各 usecase の同名メソッドをそのまま昇格させる
//     (例: 埋め込んだ *CompletePayment が CompletePayment メソッドを持つ → ラッパー側でも呼べる)
//   - 各 usecase 側で interface の signature と同名・同型のメソッドを実装する必要あり
//     (今ある Execute → method 名を usecase 名に揃えて、 input/output を publicfunctions の型に合わせる)
//
// 想定する利用 (payment.NewModule 内):
//
//	webhookCommand := paymentusecase.NewPaymentWebhookCommand(
//	    completePayment, cancelPayment, renewPayment, handleFailure,
//	)
//	return &Module{ webhookCommand: webhookCommand, ... }
type PaymentWebhookCommand struct {
	*CompletePayment
	*CancelPayment
	*RenewPayment
	*HandleFailure
}

func NewPaymentWebhookCommand(
	completePayment *CompletePayment,
	cancelPayment *CancelPayment,
	renewPayment *RenewPayment,
	handleFailure *HandleFailure,
) *PaymentWebhookCommand {
	return &PaymentWebhookCommand{
		CompletePayment: completePayment,
		CancelPayment:   cancelPayment,
		RenewPayment:    renewPayment,
		HandleFailure:   handleFailure,
	}
}
