package paymentdomain

// StripeEvent は Stripe Webhook 受信記録 + 冪等性キーを保持する。
//
// 設計 (ADR 0014):
//   - 受信した全 Stripe event を生 payload で残す (監査 / デバッグ / 再処理用)
//   - stripe_event_id (Stripe 側の evt_xxx) を UNIQUE 制約にして二重処理を物理的に防ぐ
//   - 処理完了で processed_at を SET、 失敗時は processing_error を SET
//
// migration: sql/migrations/000017_create_stripe_events.up.sql
//
// TODO: User がここから実装する
//   - field 候補:
//       id              valueobject.StripeEventID
//       stripeEventID   string  (Stripe 側の evt_xxx、 wrap するか? - 検討)
//       eventType       string  ('checkout.session.completed' 等、 文字列が多種なので enum 化しない)
//       payload         json.RawMessage  ← or map[string]any
//       processedAt     *time.Time   ← 未処理時は nil
//       processingError *string      ← 失敗時に SET
//       createdAt       time.Time
//       updatedAt       time.Time
//   - 状態遷移メソッド候補:
//       MarkProcessed(at time.Time)
//       MarkFailed(err string, at time.Time)
//   - 既に処理済みかを判定するゲッターも便利: IsProcessed() bool
