package paymentinfra

import (
	"context"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type emailTemplate struct {
	subject string
	body    string
}

// emailTemplates は outbox イベント種別 → メール本文。
// キーは valueobject.OutboxEventType の値 (relay が PublishMessage.Type に入れる)。
// もとは Lambda(outbox-consumer) にあったものを移設。
var emailTemplates = map[string]emailTemplate{
	string(valueobject.OutboxEventTypePaymentSucceeded): {
		subject: "サブスクリプションの申し込みに成功しました",
		body:    "お申し込みありがとうございます。Proプランが有効になりました。\nこのメールは送信専用です。",
	},
}

// emailPublisher は Publisher の実装で、outbox メッセージを SQS に流さず
// Mailer 経由で「その場で」メール送信する (SQS/Lambda/DynamoDB/SES を不要にする)。
type emailPublisher struct {
	mailer paymentdomain.Mailer
}

// NewEmailPublisher は Mailer でメールを直接送る Publisher を返す。
func NewEmailPublisher(mailer paymentdomain.Mailer) paymentdomain.Publisher {
	return &emailPublisher{mailer: mailer}
}

func (p *emailPublisher) Publish(ctx context.Context, msg paymentdomain.PublishMessage) error {
	tmpl, ok := emailTemplates[msg.Type]
	if !ok {
		return nil // 未対応の種別はスキップ (エラーにしない)
	}
	return p.mailer.Send(ctx, msg.Email, tmpl.subject, tmpl.body)
}
