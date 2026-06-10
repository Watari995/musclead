# ADR 0018: Webhook 処理の usecase 分割と stripe_events 保存タイミング

## ステータス
採用 (2026-06-10)

## コンテキスト

ADR 0014 で「Webhook 受信時に同期 TX で subscription 有効化」 を決定したが、 以下の詳細は決まっていなかった:

- handler / usecase の責務分割をどう設計するか
- `stripe_events` の INSERT を**いつ**、 **どの TX 内で**実行するか
- `event_type` 分岐ロジックを handler / usecase のどちらに置くか

Phase 1 実装中の議論で、 これらを確定する。

## 判断

### ① handler は parse + dispatch のみ、 各 usecase は単機能

```
[Handler]
   ↓ ParseWebhookEvent usecase (検証 + パース、 TX なし)
   ↓ event_type で分岐
   ├─ CompletePayment usecase (TX 内で stripe_events + payments + subscriptions + outbox)
   ├─ CancelPayment usecase   (同上)
   └─ RenewPayment usecase    (将来、 月次更新時)
```

- handler の責務:
  - HTTP body / header の読み取り
  - `ParseWebhookEvent` 呼び出し (署名検証)
  - `event_type` による usecase dispatch (Stripe プロトコル解釈)
  - error → HTTP status 変換
- 各 usecase の責務:
  - 単機能 (SRP)
  - TX 内で `stripe_events` 保存 + 本処理を atomic 実行

### ② `stripe_events` の INSERT は**各 usecase の TX 内**で実行

```go
// 例: CompletePayment usecase
func (uc *CompletePayment) Execute(ctx, input) error {
    return uc.txManager.Processing(ctx, func(ctx context.Context) error {
        // ⭐ TX 内で stripe_events INSERT + 本処理 を atomic に
        if err := uc.stripeEventRepo.Save(ctx, stripeEventEntity); err != nil {
            if isDuplicate(err) { return nil }  // 既に処理済み = no-op (冪等性)
            return err
        }
        // payments UPDATE / subscriptions INSERT / outbox INSERT を続けて実行
        ...
    })
}
```

### ③ `ParseWebhookEvent` は TX 外で実行

- 署名検証 (HMAC-SHA256) + JSON パースは純粋な CPU 処理、 DB 触らない
- TX 内に含めると DB lock を不必要に長く保持してしまう
- 外部 SDK 呼び出しを TX 外に出すという業界ルール

## なぜ「各 usecase の TX 内で stripe_events 保存」 なのか (核心)

### 代替案: handler が事前に saveStripeEvent usecase を呼ぶ (別 TX)

```go
// ❌ NG パターン
h.saveStripeEvent.Execute(ctx, event)    // 別 TX
h.completePayment.Execute(ctx, ...)      // 別 TX
```

### 却下理由: 整合性破綻シナリオ

1. `stripe_events` INSERT 成功 → 別 TX 確定
2. `CompletePayment` 実行中に DB エラー → TX rollback
3. **stripe_events だけ残る + payment 未確定**
4. Stripe が同じ event をリトライしてくる
5. handler が事前 `saveStripeEvent` を呼ぶ → UNIQUE 違反 (`stripe_event_id`)
6. handler は「既に処理済み」 と判断して 200 返す
7. **CompletePayment が永久に実行されず、 payment が永遠に未確定状態**

### 採用案: 同一 TX 内で atomic

- `stripe_events` INSERT + 本処理が常に atomic
- どちらかが失敗すれば両方ロールバック → 次のリトライで再実行可能
- **「stripe_events の存在 = 本処理完了」 が DB レベルで保証される**

## 実装上のルール

| 処理 | TX 内 / 外 | 場所 |
|---|---|---|
| `ParseWebhookEvent` (署名検証 + パース) | **TX 外** | usecase (薄い wrapper) |
| `stripe_events` INSERT | **TX 内** | 各 event_type の usecase 内 |
| `payments` UPDATE | TX 内 | 同上 |
| `subscription_orders` UPDATE | TX 内 | 同上 |
| `subscriptions` INSERT | TX 内 | 同上 |
| `outbox_events` INSERT | TX 内 | 同上 |
| 即時 SNS publish | **TX 外** | usecase が TX commit 後に実行 (ADR 0015) |

## usecase 一覧 (Phase 1 で実装)

| usecase | 引数 | 責務 |
|---|---|---|
| `ParseWebhookEvent` | payload + sigHeader | StripeClient の薄い wrapper、 handler から呼ばれる |
| `InitiatePayment` | user_id, price_id 等 | 申込開始 (Stripe Checkout Session 作成) |
| `CompletePayment` | StripeEventID + EventType + Payload + ... | success Webhook 時の処理 |
| `CancelPayment` | 同上 | cancel Webhook 時の処理 |
| `RenewPayment` (将来) | 同上 | 月次更新 Webhook 時の処理 |
| `GetPayment` (Query) | payment_id | 単体取得 |

## 関連 ADR

- [ADR 0013](0013-purchase-payment-separation.md): purchase / payment 分離
- [ADR 0014](0014-webhook-idempotency-and-retry.md): Webhook 冪等性
- [ADR 0015](0015-outbox-pattern-and-async-mail.md): Outbox パターン (TX 外 SNS publish)

## 影響

### Phase 1 実装に直接反映
- `internal/payment/internal/handler/webhook_handler.go`: parse + dispatch
- `internal/payment/internal/usecase/parse_webhook_event.go`: StripeClient 薄い wrapper
- `internal/payment/internal/usecase/complete_payment.go`: TX 内 atomic
- `internal/payment/internal/usecase/cancel_payment.go`: TX 内 atomic

### handler が知っていい範囲
- HTTP プロトコル (body / header / status)
- Stripe 固有の event_type 文字列 ('checkout.session.completed' 等) ← これは Stripe Webhook プロトコルそのものなので handler に置いて OK
