package paymentinfra

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
)

// MVPでは外部サービスのリトライに任せる
type ExternalRetryStrategy struct{}

// errをそのまま返すことで500のエラーを伝播させる
func (s *ExternalRetryStrategy) OnFailure(ctx context.Context, event *paymentdomain.StripeEvent, err error) error {
	return err
}
