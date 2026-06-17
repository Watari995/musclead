package paymentdomain

import "context"

// PublishMessage は outbox イベントを外部 (SQS) に流す時の自己完結メッセージ。
// Lambda は VPC 外 = DB に触れないため、 メール送信に必要な情報をここに全部詰める (ADR 0020 ③)。
type PublishMessage struct {
	EventID string // 冪等キー (outbox event id)。 consumer が重複排除に使う
	Type    string // イベント種別 (例: payment_succeeded)。 consumer がテンプレを選ぶ
	Email   string // 宛先 (relay が enrich)
}

// Publisher は outbox メッセージの送信先 (SQS) を抽象化する port。
// usecase / domain は AWS SDK を直接見ない (StripeClient と同じ ACL)。
type Publisher interface {
	Publish(ctx context.Context, msg PublishMessage) error
}
