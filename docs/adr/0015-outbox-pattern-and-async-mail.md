# ADR 0015: Outbox パターンとメール送信の非同期化 (SNS + SQS + Lambda)

## ステータス
採用 (2026-06-09)

## コンテキスト

[ADR 0014](0014-webhook-idempotency-and-retry.md) で「subscription 有効化は Webhook 内で同期実行」 と決めた。 一方、 **メール送信は処理が重く、 失敗時のリトライが必要**なため非同期化する。

求める性質:
1. メール送信失敗が決済処理本体を巻き込まない
2. 失敗時に自動リトライ + DLQ 退避
3. データロストしない (outbox で持続性を保証)
4. コストを抑える (個人開発、 worker 常時稼働は避ける)

## 判断

### ① Outbox パターン + SNS + SQS + Lambda 構成

```
[Webhook 受信、 TX 内]
  ├─ payments UPDATE
  ├─ subscriptions INSERT
  └─ outbox_events INSERT  ⭐ atomic に「通知すべきイベント」 を記録
       │
[TX 外]
  └─ 即時 SNS publish 試行
       │ 成功 → outbox_events.published_at = NOW()
       │ 失敗 → outbox_events.published_at = NULL のまま残置
       ▼
  [SNS Topic: payment-events]
       │
       ▼
  [SQS Queue: email-events] (+ DLQ)
       │
       ▼
  [Lambda: email worker] → SES SendEmail → emails INSERT
```

### ② Outbox INSERT は **TX 内**で行う (核心)

```go
txRunner.RunInTx(func(tx) {
    paymentRepo.Save(tx, payment)
    subscriptionRepo.Save(tx, subscription)
    outboxRepo.Save(tx, outboxEvent)   // ⭐ 同じ TX
})
```

これにより:
- TX commit 成功 → outbox に記録あり → 必ずいつかは publish される
- TX rollback → outbox に記録なし → イベントは発生しなかったことに

→ 「DB と外部メッセージブローカー」 の 2-phase commit 問題を回避。

### ③ 即時 publish + outbox-relay の二段構え (failsafe)

通常時:
- Webhook handler が TX commit 後すぐ `sns.Publish()` を呼ぶ → 200ms 以内に SQS に届く
- 成功時は `outbox_events.published_at = NOW()` を立てる

異常時 (SNS 一時障害、 ネットワーク切断):
- `published_at = NULL` のまま outbox に残る
- 1 分間隔の `outbox-relay` Lambda が拾って再送

```go
// outbox-relay Lambda (EventBridge schedule で 1 分ごと起動)
events := outboxRepo.FindPending(ctx, 100)   // WHERE published_at IS NULL
for _, e := range events {
    sns.Publish(e)
    outboxRepo.MarkPublished(ctx, e.ID)
}
```

→ 通常 1 秒、 異常時でも最大 1 分でメール送信が起動される。

### ④ inbox_events は作らない (subscription 同期処理のため不要)

当初設計では「purchase 側で `inbox_events` テーブルを作り、 SQS 受信時の冪等性を担保」 を予定。 ただし [ADR 0014](0014-webhook-idempotency-and-retry.md) で subscription 有効化を同期化したため、 purchase 側に SQS consumer が存在しない → `inbox_events` 不要。

メール送信用 (Lambda email-worker) の冪等性は `emails.source_event_id` UNIQUE で代替:

```sql
CREATE TABLE emails (
  id              BINARY(16)   NOT NULL,
  user_id         BINARY(16)   NOT NULL,
  template_type   VARCHAR(50)  NOT NULL,
  to_address      VARCHAR(255) NOT NULL,
  subject         VARCHAR(255) NOT NULL,
  source_event_id BINARY(16)   NOT NULL,        -- ⭐ SQS message の冪等性キー
  sent_at         DATETIME(6)  NULL,
  failed_at       DATETIME(6)  NULL,
  failure_reason  TEXT         NULL,
  ses_message_id  VARCHAR(255) NULL,
  created_at      DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at      DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  UNIQUE KEY uk_source (source_event_id)        -- 二重送信防止
);
```

### ⑤ Lambda + EventBridge 構成 (Fargate worker は不採用)

| 構成 | 月額 (常時稼働) | 評価 |
|---|---|---|
| Fargate worker (常時稼働 × 2 タスク) | $18 | 個人開発には重い |
| Fargate worker (`desired_count=0` で必要時起動) | $0-2 | 本番運用には不向き (申込時に worker 起動が必要) |
| **Lambda + EventBridge** ✅ | **$0** (free tier 内) | 本番常時運用可、 コスト最小 |

Lambda 構成:
- `outbox-relay`: EventBridge 1 分間隔 schedule、 Reserved Concurrency = 1 (同時実行防止)
- `email-worker`: SQS trigger、 失敗時は自動リトライ + DLQ
- 実装形式: Container image (ECR から pull、 既存 Docker pipeline 流用)

### ⑥ メール送信は Stripe デフォルトメールに任せず自前実装

ユーザー通知の主導権を持ちたいため自前 SES 送信。 4 種類のテンプレートを送信:

| トリガー | テンプレート | 内容 |
|---|---|---|
| PaymentSucceeded | `payment_succeeded` | 「Pro へようこそ、 X 月 X 日まで有効」 |
| PaymentFailed | `payment_failed` | 「決済に失敗、 再度お試しください」 |
| PaymentPastDue (Dunning) | `dunning` | 「カード更新が必要 + Customer Portal リンク」 |
| PaymentCanceled | `canceled` | 「解約受付、 X 月 X 日まで Pro 利用可」 |

## なぜ Outbox パターンを採用したか

代替案: TX commit 後に直接 SQS send。

却下理由 (2-phase commit 問題):
- TX commit 成功 + SQS send 失敗 → DB は更新済みだが下流に通知されない (永久ロスト)
- TX commit 失敗の後に SQS send 成功 → 通知だけ流れて DB は古いまま (二重処理の元)

Outbox パターンなら **1 つの DB TX 内**で「アプリ状態更新 + イベント記録」 を atomic に終わらせ、 後から非同期で配送できる。

## なぜ SNS を経由するか (直接 SQS でない)

代替案: outbox-relay が直接 SQS send。

却下理由:
- 将来 subscriber を増やす (例: Slack 通知、 分析基盤への送信) ときに、 SQS を増やすだけで対応したい
- SNS fan-out が標準パターン
- SNS / SQS のコストは無視できるレベル

## なぜ purchase 側に inbox_events を作らないか

過去のドラフトでは作る予定だったが、 [ADR 0014](0014-webhook-idempotency-and-retry.md) で subscription 有効化を同期化したため不要に。

| 設計 | inbox_events |
|---|---|
| 当初: 非同期 subscription 有効化 (purchase worker が SQS から受信) | ⭕ 必要 |
| 現在: 同期 subscription 有効化 (Webhook 内 TX) | ✗ 不要 |

メール送信側 (Lambda email-worker) の冪等性は emails テーブル UNIQUE で対応。

## やらないこと

- 自前メッセージブローカー (SNS / SQS で十分)
- WebSocket / SSE での結果プッシュ (success ページのポーリングで十分)
- DLQ メッセージの自動再投入 (手動で対応、 個人開発)
- Slack / CloudWatch アラート (本格運用フェーズで)

## 影響

### 新規テーブル
- `outbox_events` (payment)
- `emails` (email context)

### 新規 Lambda
- `outbox-relay`: EventBridge schedule (1 min)、 Container image
- `email-worker`: SQS trigger、 Container image

### 新規 module
- `internal/email/` (template 管理、 送信履歴)

### Terraform
- SNS Topic (payment-events)
- SQS Queue (email-events) + DLQ
- Lambda function × 2 + IAM role
- EventBridge schedule

### コスト
- Lambda + SNS + SQS + EventBridge: 月 **$0.04 程度** (free tier 内)

## 関連 ADR

- [ADR 0013](0013-purchase-payment-separation.md): purchase / payment 分離
- [ADR 0014](0014-webhook-idempotency-and-retry.md): Webhook 冪等性
