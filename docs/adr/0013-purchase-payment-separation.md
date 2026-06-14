# ADR 0013: 購入 (purchase) と決済 (payment) の bounded context 分離

## ステータス
採用 (2026-06-09)

## コンテキスト

Pro サブスクリプションを実装するにあたり、 「**申込**」 と「**決済**」 を別の関心事として扱う必要がある。

| 関心事 | 主な責務 |
|---|---|
| 購入 (purchase) | 「ユーザーが何を買おうとしたか」 のライフサイクル、 申込・成功・失敗、 サブスクの権利状態 |
| 決済 (payment) | 「実際の支払い処理」、 Stripe との通信、 状態遷移、 監査ログ |

SODA の事例 ([参考記事](https://zenn.dev/team_soda/articles/16b742f8a58c45)) でも 2 つの bounded context を**明確に分離**しており、 「購入は決済のことを知るのは OK、 逆は不可」 と一方向依存を厳格に守っている。 この方針を musclead でも踏襲する。

## 判断

### ① bounded context を 2 つに分離 (purchase / payment)

```
internal/
  purchase/      ← 申込のオーケストレータ (薄い、 ロジック少)
  payment/       ← Stripe との統合 (厚い、 状態管理が主)
```

### ② 依存方向: `purchase → payment` の一方向

- `purchase` は `payment` の公開 interface に依存して OK
- `payment` は `purchase` を一切知らない (= `payment` 配下に `purchase` への参照を持たない)
- 逆方向の通信が必要な場合 (例: Webhook 受信時に subscription を有効化する) は、 **Webhook handler が両 context を呼ぶオーケストレーター**として扱う ([ADR 0014](0014-webhook-idempotency-and-retry.md) 参照)

### ③ Facade パターン: payment の公開 API

SODA 記事に倣い、 payment は **4 つのメソッド**を購入側に公開する:

```go
// payment/payment.go
type Module struct {
    Command PaymentCommand
    Query   PaymentQuery
}

type PaymentCommand interface {
    InitiatePayment(ctx context.Context, input InitiatePaymentInput) (*InitiatePaymentOutput, error)
    CompletePayment(ctx context.Context, input CompletePaymentInput) error
    CancelPayment(ctx context.Context, input CancelPaymentInput) error
    CapturePayment(ctx context.Context, input CapturePaymentInput) error  // サブスクでは未使用、 将来用
}
```

既存 musclead の `UserCommand` / `WeightCommand` 等のパターンを踏襲。

### ④ 集約間の参照: ID のみ (DDD 原則)

- `subscription_orders.payment_id` を保持 (purchase → payment の参照は ID で)
- `payments` 側は `order_id` を持たない (依存方向違反になるため)
- 詳細データが必要なら `payment.Query.GetPayment(payment_id)` 経由で取得

### ⑤ 抽象化レイヤー: 将来の決済 SaaS 追加に備える

Stripe 固有の処理は `StripeWebhookProcessor` / `StripeClient` 内に閉じ、 上位 (handler / usecase) は interface 越しに利用。

```go
// 将来 PAY.JP 追加時:
type WebhookProcessor interface {
    VerifySignature(req) (Event, error)
    ProcessEvent(ctx, event) error
}
type StripeWebhookProcessor struct { ... }
type PayJPWebhookProcessor struct { ... }   // 追加するだけ
```

Composition Root で DI 注入を切替えるだけで決済 SaaS の追加が可能。

### ⑥ 主要テーブルの配置

| テーブル | context | 役割 |
|---|---|---|
| `subscription_orders` | purchase | 申込トリガー (1 回限り、 履歴) |
| `subscriptions` | purchase | 継続的な権利状態 (active / canceled / expired) |
| `payments` | payment | Stripe 決済本体 (1 record = 1 Stripe Subscription) |
| `payment_events` | payment | 監査ログ (append-only、 状態遷移記録) |
| `stripe_events` | payment | Stripe Webhook 受信記録 (生 payload + 冪等性) |
| `outbox_events` | payment | メール送信用 outbox ([ADR 0015](0015-outbox-pattern-and-async-mail.md)) |
| `idempotency_keys` | shared | API 重複排除 middleware ([ADR 0014](0014-webhook-idempotency-and-retry.md)) |
| `emails` | email (新規) | メール送信履歴 + 冪等性 |

### ⑦ 「申込トリガー」 と「権利状態」 を分離 (subscription_orders + subscriptions)

DDD 的に重要なポイント: 1 つのテーブルで両方を扱うと state machine が複雑化し、 再加入の表現も不自然になる。

| 集約 | 何を表現するか |
|---|---|
| `subscription_orders` | **申込意思** (1 回限りのイベント、 pending → succeeded / failed) |
| `subscriptions` | **権利状態** (continuous、 active → canceled → expired) |

Stripe 自体も `Checkout Session` (申込) と `Subscription` (権利) を別エンティティで管理しているのと一致。

## なぜ「決済 / 購入の分離」 を採用したか

代替案: `billing` 1 つの context にまとめる。

却下理由:
- 「申込履歴」 と「決済状態」 と「権利状態」 が混ざり、 集約境界が曖昧になる
- 将来複数の決済 SaaS を追加する際、 関心の分離がないと拡張が破綻する
- SODA 流儀 (本 ADR の参考記事) を踏襲することで、 入社時に即戦力で議論に参加できる

## なぜ Customer 管理に専用テーブルを作らないか

代替案: `payment_customers` テーブルで `user_id` UNIQUE で Stripe Customer ID を管理。

却下理由:
- `payments.stripe_customer_id` を SELECT して再利用するだけで実用上十分
- 同時並行リクエストによる二重作成は **idempotency middleware で防げる** (クライアントから同じ Idempotency-Key で来た場合)
- テーブル増加コストに見合う価値が低い

ただし、 「同一 user に必ず 1 つの Stripe Customer」 を強制したい時 (例: 法人プラン追加時) は専用テーブル化を再検討。

## やらないこと (本 ADR スコープ外)

- 複数決済 SaaS の同時稼働 (本 ADR では「追加可能な構造を残す」 だけ)
- 集約間の同期手段 (Webhook 内 TX or イベント駆動) の詳細 → [ADR 0014](0014-webhook-idempotency-and-retry.md)
- メール送信の非同期構成 → [ADR 0015](0015-outbox-pattern-and-async-mail.md)

## 影響

### 新規 module
- `internal/purchase/` (申込ライフサイクル + サブスク権利)
- `internal/payment/` (Stripe 統合 + 監査)
- `internal/email/` (SES 送信、 [ADR 0015](0015-outbox-pattern-and-async-mail.md) で詳細)

### 既存 module への影響
- `user/` は変更なし
- ただし `/users/me` で subscription 情報を返すため、 user handler が `purchase.SubscriptionQuery` に依存する (handler レイヤーでの統合は OK)

### Composition Root (main.go)
- `paymentModule := payment.NewModule(...)`
- `purchaseModule := purchase.NewModule(db, paymentModule.Command)` ← Command 注入
- これで `purchase → payment` の一方向依存が型レベルで保証される

## 関連 ADR

- [ADR 0002](0002-ddd-modular-monolith.md): DDD + Modular Monolith の全体方針
- [ADR 0004](0004-module-public-interface.md): モジュール公開 interface (Command / Query パターン)
- [ADR 0012](0012-premium-features-overview.md): プレミアム機能の全体方針
- [ADR 0014](0014-webhook-idempotency-and-retry.md): Webhook 冪等性とリトライ戦略
- [ADR 0019](0019-billing-module-webhook-orchestrator.md): Webhook handler は payment ではなく新規 `billing` module に置く (本 ADR の「handler 例外」 を構造で具現化)

## 更新履歴

- 2026-06-11: Webhook handler の置き場所を [ADR 0019](0019-billing-module-webhook-orchestrator.md) で `internal/billing/` に変更。 本 ADR の依存方向 (`purchase → payment` 一方向) は維持。
