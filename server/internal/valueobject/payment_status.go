package valueobject

import "errors"

type PaymentStatusCode string

const (
	PaymentStatusPending   PaymentStatusCode = "pending"
	PaymentStatusSucceeded PaymentStatusCode = "succeeded"
	PaymentStatusFailed    PaymentStatusCode = "failed"
	PaymentStatusCanceled  PaymentStatusCode = "canceled"
)

var ErrInvalidPaymentStatus = errors.New("invalid payment status")

type PaymentStatus struct {
	LiteralBase[string]
}

// 外部からの文字列からPaymentStatusを作成する
func NewPaymentStatusFromString(s string) (*PaymentStatus, error) {
	switch PaymentStatusCode(s) {
	case PaymentStatusPending, PaymentStatusSucceeded, PaymentStatusFailed, PaymentStatusCanceled:
		return &PaymentStatus{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidPaymentStatus
	}
}

// CodeからPaymentStatusを作成する
func NewPaymentStatusFromCode(c PaymentStatusCode) PaymentStatus {
	return PaymentStatus{LiteralBase: LiteralBase[string]{v: string(c)}}
}
