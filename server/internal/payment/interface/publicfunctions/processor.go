// Package publicfunctions の Stripe Webhook 検証 / パース interface。
//
// 設計 (ADR 0019): billing module が Stripe 固有の署名検証を payment に閉じたまま
// 呼べるように、 publicfunctions 経由で公開する。 billing から見ると「Stripe SDK を
// 直接触らずに Stripe Event を取り出せるラッパー」 として扱える。
package publicfunctions

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

// StripeEvent は署名検証 + パース後に handler が受け取る正規化された event。
type StripeEvent struct {
	StripeEventID string               // Stripe 発行の一意 ID (evt_xxx)
	EventType     string               // 'checkout.session.completed' 等
	Payload       valueobject.Metadata // event.Data の正規化済み中身
}

// StripeWebhookProcessor は HTTP body / Stripe-Signature を受け取り、 検証済み event を返す。
//
// 設計メモ (ADR 0018 ③): TX 外で実行 (純粋 CPU 処理、 DB を触らない)。
type StripeWebhookProcessor interface {
	ParseAndVerify(ctx context.Context, payload []byte, signatureHeader string) (StripeEvent, error)
}
