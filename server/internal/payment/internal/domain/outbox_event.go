package paymentdomain

// OutboxEvent は payment context が発行する通知用 outbox (ADR 0015)。
//
// 設計:
//   - Webhook 処理の TX 内で INSERT し、 TX 外で SNS publish する
//   - 即時 publish 成功時は published_at を SET
//   - 即時 publish 失敗時は 1 分後の outbox-relay Lambda が拾う (failsafe)
//
// migration: sql/migrations/000018_create_outbox_events.up.sql
//
// TODO: User がここから実装する
//   - field 候補:
//       id           valueobject.OutboxEventID
//       eventType    valueobject.OutboxEventType (要 VO 化、 'PaymentSucceeded' / 'PaymentFailed' / 'PaymentCanceled' / 'PaymentRenewed')
//       aggregateID  valueobject.PaymentID  (発生源、 polling worker のグルーピング用)
//       payload      json.RawMessage  ← SQS に流す本文
//       publishedAt  *time.Time       ← 未配信時は nil
//       publishError *string          ← 失敗時に SET
//       createdAt    time.Time
//       updatedAt    time.Time
//   - 状態遷移メソッド候補:
//       MarkPublished(at time.Time)
//       MarkPublishFailed(err string)
//   - 既に publish 済みかのゲッター: IsPublished() bool
