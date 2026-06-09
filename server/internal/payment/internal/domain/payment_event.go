package paymentdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// PaymentEvent は payment 集約の状態遷移を時系列で記録する **append-only** な監査ログ。
//
// 設計 (ADR 0014):
//   - 状態遷移 (initiated / succeeded / failed / canceled / renewed) を時系列で記録
//   - updated_at なし: 一度書いたら書き換えない (append-only)
//   - 将来の決済 SaaS 追加 (PAY.JP 等) に備え、 決済サービス非依存の抽象表現
//
// migration: sql/migrations/000016_create_payment_events.up.sql
//
// TODO: User がここから実装する
//   - Payment と同じパターン (private fields + getters + Create/New)
//   - field 候補:
//       id           valueobject.PaymentEventID
//       paymentID    valueobject.PaymentID
//       eventType    valueobject.PaymentEventType (要 VO 化)
//       metadata     json.RawMessage  ← or map[string]any (DB 側は JSON)
//       createdAt    time.Time
//   - Update メソッド不要 (append-only なので状態変更しない)
//   - CreatePaymentEvent: 新規作成、 ID と createdAt は内部で生成
//   - NewPaymentEvent: 全 field 受け取り (repo から復元用)

type PaymentEvent struct {
	id        valueobject.PaymentEventID
	paymentID valueobject.PaymentID
	eventType valueobject.PaymentEventType
	metadata  valueobject.Metadata
	createdAt time.Time
}

func (p *PaymentEvent) ID() valueobject.PaymentEventID {
	return p.id
}

func (p *PaymentEvent) PaymentID() valueobject.PaymentID {
	return p.paymentID
}

func (p *PaymentEvent) EventType() valueobject.PaymentEventType {
	return p.eventType
}

func (p *PaymentEvent) Metadata() valueobject.Metadata {
	return p.metadata
}

func (p *PaymentEvent) CreatedAt() time.Time {
	return p.createdAt
}

func CreatePaymentEvent(
	paymentID valueobject.PaymentID,
	eventType valueobject.PaymentEventType,
	metadata valueobject.Metadata,
) *PaymentEvent {
	return &PaymentEvent{
		id:        valueobject.NewPrimaryID[valueobject.PaymentEventID](),
		paymentID: paymentID,
		eventType: eventType,
		metadata:  metadata,
		createdAt: time.Now(),
	}
}

func NewPaymentEvent(
	id valueobject.PaymentEventID,
	paymentID valueobject.PaymentID,
	eventType valueobject.PaymentEventType,
	metadata valueobject.Metadata,
	createdAt time.Time,
) *PaymentEvent {
	return &PaymentEvent{
		id:        id,
		paymentID: paymentID,
		eventType: eventType,
		metadata:  metadata,
		createdAt: createdAt,
	}
}
