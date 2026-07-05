// Package paymentdomain は payment 集約の domain 層を定義する。
//
// 含むもの:
//   - entity: Payment, PaymentEvent, StripeEvent
//   - VO: PaymentStatus, PaymentID 等 (大半は internal/valueobject に置く)
//   - interface:
//       - PaymentRepository / PaymentEventRepository / StripeEventRepository
//       - StripeClient (Stripe SDK を抽象化)
//       - RetryStrategy (ADR 0014 で抽象化、 将来 PAY.JP 等を追加できるように)
//
// OutboxEvent / OutboxEventRepository は internal/shared/domain に移設済み
// (複数モジュールで outbox パターンを共有するため)。
//
// 設計参考:
//   - ADR 0013 (purchase / payment 分離)
//   - ADR 0014 (Webhook 同期処理 + 冪等性 3 レイヤー)
//   - 既存 internal/weight/internal/domain/ の entity / repository パターン
package paymentdomain
