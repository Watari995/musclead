# ADR 0014: Webhook の冪等性と即時 subscription 有効化、 リトライ戦略

## ステータス
採用 (2026-06-09)

## コンテキスト

Stripe Webhook を受信した時、 以下を達成する必要がある:

1. **即時 Pro 化** (ユーザー体感 1 秒以内に subscription が active)
2. **冪等性** (Stripe の自動リトライによる重複受信、 クライアント連打、 Lambda 二重実行を全て吸収)
3. **失敗時の安全な回復** (一時障害は自動回復、 永続障害は手動再送可能)

これらは同時に成立させる必要があり、 設計上の鍵となる。

## 判断

### ① subscription 有効化は Webhook 内で同期実行 (TX 内)

```go
// Webhook handler
func (h *WebhookHandler) Handle(req) {
    // 1. stripe_events INSERT (冪等性、 UNIQUE 制約で重複弾く)
    err := stripeEventsRepo.Save(ctx, event)
    if errors.Is(err, ErrDuplicate) { return 200 }

    // 2. TX 内で全部 atomic に実行
    err := txRunner.RunInTx(ctx, func(tx Tx) error {
        h.paymentCommand.CompletePayment(ctx, tx, event)        // payment 集約
        h.purchaseCommand.ActivateSubscription(ctx, tx, ...)    // ⭐ subscription INSERT
        h.outboxRepo.Save(ctx, tx, emailEvent)                  // メール送信用 outbox
        return nil
    })

    // 3. TX 外: メール用 SNS publish (即時、 失敗時は outbox-relay が拾う)
    h.snsClient.Publish(emailEvent)

    return 200
}
```

| 処理 | 同期 / 非同期 | 理由 |
|---|---|---|
| stripe_events INSERT | 同期 (TX 内) | 冪等性キー |
| payments UPDATE | 同期 (TX 内) | Pro 化の前提条件 |
| subscription_orders UPDATE | 同期 (TX 内) | 申込ライフサイクル完結 |
| **subscriptions INSERT** | **同期 (TX 内)** | **即時 Pro 化のため** ⭐ |
| outbox INSERT | 同期 (TX 内) | atomic 保証 |
| SNS publish (メール用) | 非同期 (TX 外) | メール送信は時間かかる、 失敗時リトライ可 |

→ 通常時 **100ms 以内に subscriptions が active**。

### ② 冪等性は 3 レイヤーで担保

| レイヤー | テーブル | 役割 |
|---|---|---|
| 1. Stripe Webhook 重複 | `stripe_events` (`stripe_event_id` UNIQUE) | Stripe 自動リトライによる重複受信を弾く |
| 2. クライアント API 重複 | `idempotency_keys` (middleware) | ユーザーの連打 / ネットワークリトライ |
| 3. メール送信 (Lambda 二重実行) | `emails` (`source_event_id` UNIQUE) | SQS at-least-once 配送 |

### ③ subscriptions INSERT も冪等に作る

Stripe リトライで Webhook が 2 回来た時、 subscription が二重 INSERT されないように。

```go
func (uc *ActivateSubscription) Execute(ctx, tx, paymentID) error {
    existing, _ := uc.subscriptionRepo.FindByPaymentID(ctx, tx, paymentID)
    if existing != nil { return nil }  // ⭐ 既に有効化済み、 何もしない

    sub := domain.NewSubscription(...)
    return uc.subscriptionRepo.Save(ctx, tx, sub)
}
```

### ④ クライアント API 重複排除: idempotency_keys + middleware

クライアントは `Idempotency-Key: <UUID>` ヘッダを送る。 サーバはレスポンス本体 (`response_snapshot`) を保存し、 同 key 再送時は同じレスポンスを返す。 Stripe API と同じパターン。

#### スキーマ

```sql
CREATE TABLE idempotency_keys (
  id              BINARY(16)   NOT NULL,
  idempotency_key VARCHAR(255) NOT NULL,
  user_id         BINARY(16)   NOT NULL,         -- 他ユーザーの key 偽装防止
  request_path    VARCHAR(255) NOT NULL,         -- 違う API への同 key 検知
  request_hash    VARCHAR(64)  NOT NULL,         -- body 改ざん検知
  response_status SMALLINT     NULL,
  response_body   JSON         NULL,
  completed_at    DATETIME(6)  NULL,             -- NULL = 処理中
  created_at      DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at      DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_key (user_id, idempotency_key)
);
```

middleware は全 POST/PUT/PATCH に適用 (将来の cancel / refund / plan_change も同じ仕組みで保護される)。

### ⑤ Webhook 失敗時のリトライ: 各決済 SaaS の自動リトライに任せる (MVP)

通常時:
- TX 失敗 → 500 返却 → Stripe が自動リトライ (指数バックオフ、 最大 3 日間)
- 全処理が冪等なので 2 回目以降の処理も安全

| 失敗箇所 | 対応 |
|---|---|
| stripe_events INSERT 失敗 | 500 → Stripe リトライ |
| TX 失敗 | 500 → Stripe リトライ (subscription / payment / outbox が atomic にロールバック) |
| 3 日経過後の失敗 | Stripe Dashboard で手動 resend |
| 永続的なバグ | バグ修正後、 Stripe Dashboard で手動 resend |

### ⑥ 将来の拡張: RetryStrategy interface で切替可能に

複数決済 SaaS を採用する将来に向けて、 リトライ戦略を抽象化:

```go
type RetryStrategy interface {
    OnFailure(ctx context.Context, event Event, err error) error
}

// 現在 (MVP): 外部サービスのリトライに任せる
type ExternalRetryStrategy struct{}
func (s *ExternalRetryStrategy) OnFailure(ctx, event, err) error {
    return err   // 5xx を返してリトライしてもらう
}

// 将来 (複数 SaaS 採用時): 自前で failed_webhook_events に記録し、 自前 Lambda で再処理
type SelfManagedRetryStrategy struct {
    failedEventRepo FailedEventRepository
}
func (s *SelfManagedRetryStrategy) OnFailure(ctx, event, err) error {
    return s.failedEventRepo.Save(ctx, event, err)  // 200 を返して自前リトライへ
}
```

切替時のコスト: Composition Root で 1 行差し替え + 必要なら `failed_webhook_events` テーブル migration。

## なぜ subscription 有効化を非同期にしないか

最初は「`purchase → payment` の依存方向を守るため event 駆動 (SNS + SQS + Lambda) で subscription を有効化」 を検討したが却下。

| 観点 | 同期 (採用) | 非同期 (却下) |
|---|---|---|
| Pro 化遅延 | < 100ms | 1 秒〜 1 分 |
| 失敗時のロスト | TX で atomic、 不可能 | outbox failsafe 必要 |
| 依存方向の維持 | handler が両 context を呼ぶ orchestrator で OK | event 駆動でも OK |
| 実装複雑度 | 低 | 高 |

「即時 Pro 化」 という UX 要件が強く、 同期処理が正解。 handler レイヤーは複数 context を統合する権限を持つので依存方向違反ではない。

## なぜ「Stripe リトライ任せ」 で十分か

Stripe の Webhook リトライ仕様:
- 5xx / timeout で自動リトライ
- 指数バックオフ (5 秒、 5 分、 30 分、 2 時間、 5 時間、 ...)
- **最大 3 日間継続**
- 最終失敗時は Stripe Dashboard で確認、 手動 resend 可能

3 日のリトライ枠は実用的に十分。 短期障害 (数時間のサーバダウン、 DB 一時障害) は完全自動回復。

長期障害 (3 日以上のダウン、 永続バグ) は手動運用必要だが、 個人開発レベルでは許容。

## やらないこと

- 自前リトライ機構の MVP 実装 (将来 `SelfManagedRetryStrategy` で対応)
- `failed_webhook_events` テーブル (上記と同じく将来)
- CloudWatch アラート / Slack 通知 (本格運用フェーズで追加)
- Webhook 失敗の自動再処理バッチ (Stripe Dashboard で手動再送で十分)

## 影響

### 新規テーブル
- `idempotency_keys` (shared)
- `stripe_events` (payment)

### 新規ファイル
- `shared/middleware/idempotency.go`
- `payment/internal/handler/webhook_handler.go`
- `payment/internal/usecase/handle_webhook.go`
- `payment/internal/usecase/complete_payment.go`
- `payment/internal/domain/retry_strategy.go` (interface + ExternalRetryStrategy 実装)

### 既存ファイル変更
- `cmd/server/main.go`: idempotency middleware の全 API 適用、 RetryStrategy DI 注入

## 関連 ADR

- [ADR 0013](0013-purchase-payment-separation.md): purchase / payment 分離
- [ADR 0015](0015-outbox-pattern-and-async-mail.md): Outbox + メール非同期
- [ADR 0017](0017-stripe-integration-details.md): Stripe 統合の細部
- [ADR 0019](0019-billing-module-webhook-orchestrator.md): Webhook handler の置き場所を `internal/billing/` に変更

## 更新履歴

- 2026-06-11: Webhook handler の置き場所を [ADR 0019](0019-billing-module-webhook-orchestrator.md) で `internal/billing/` に移設。 「handler は両 context のオーケストレーター」 の趣旨は維持し、 ファイル位置で構造的に表現する形に変更。 単一 TX で全 atomic の図 (本 ADR ①) は MVP では「各 usecase 独自 TX + 冪等性で結果整合」 にダウングレード ([ADR 0019](0019-billing-module-webhook-orchestrator.md) ⑥)。
