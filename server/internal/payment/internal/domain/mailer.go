package paymentdomain

import "context"

// Mailer はメール送信を抽象化する port。実装は infra (Resend 等)。
// Publisher と同様、usecase / domain は送信手段の詳細 (AWS/HTTP) を見ない。
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}
