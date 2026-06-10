package paymentusecase

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// GetPaymentInput は単体取得の入力。
type GetPaymentInput struct {
	PaymentID valueobject.PaymentID
}

// GetPaymentOutput は単体取得の結果 (DTO)。
// public 公開時は dto.PaymentDTO 等に変換するが、 内部 usecase は entity をそのまま返す。
type GetPaymentOutput struct {
	Payment *paymentdomain.Payment
}

// GetPayment は payment 単体を取得する Query 系 usecase。
//
// 用途:
//   - purchase 集約から checkout_url のリトライ取得 (ADR 0017)
//   - 他 context が /payment/{id} 相当の API を呼ぶ時の base
type GetPayment struct {
	paymentRepo paymentdomain.PaymentRepository
}

func NewGetPayment(paymentRepo paymentdomain.PaymentRepository) *GetPayment {
	return &GetPayment{paymentRepo: paymentRepo}
}

// Execute は payment_id で取得する。 見つからない時は (output{nil}, nil)。
//
// TODO (User 実装):
//
//	payment, err := uc.paymentRepo.FindByID(ctx, input.PaymentID)
//	if err != nil { return GetPaymentOutput{}, err }
//	return GetPaymentOutput{Payment: payment}, nil
func (uc *GetPayment) Execute(ctx context.Context, input GetPaymentInput) (GetPaymentOutput, error) {
	return GetPaymentOutput{}, errNotImplemented
}
