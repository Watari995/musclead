# ADR 0019: Stripe Webhook の受信を `billing` module に分離 (新規オーケストレーター層)

## ステータス
採用 (2026-06-11)

## コンテキスト

Phase 1 では Stripe Webhook 受信処理を `internal/payment/internal/handler/webhook_handler.go` に実装した。 当時の認識:

- ADR 0014 ①: 「Webhook 受信時に subscription 有効化を同期 TX 実行」
- ADR 0014 末尾: 「handler レイヤーは複数 context を統合する権限を持つので依存方向違反ではない」

しかし Phase 2 で purchase 集約を立てて実装してみると、 以下の構造的問題が顕在化した:

1. **payment が purchase の publicfunctions を import する**ことになり、 ADR 0013 ② の言明
   > `payment` は `purchase` を一切知らない (= `payment` 配下に `purchase` への参照を持たない)

   と矛盾する。 ADR 0014 の「handler 例外」 は文面上認めているが、 ファイル位置が `payment/internal/handler/` だと「payment が purchase を知っている」 ように見え、 依存方向の表明としてノイズが大きい。

2. **Composition Root で循環依存が発生する**:
   - purchase は paymentCommand を受けて構築する
   - payment の webhook handler は purchaseCommand を呼びたい
   - 組み立て順を決められない

3. Webhook handler の本質は「Stripe Webhook → アプリ内 action」 の**アダプター** (Hexagonal Architecture でいう Driving Adapter / Inbound) であり、 純粋には payment にも purchase にも属さない。 ベンダー固有の HTTP プロトコル解釈と、 自アプリの bounded context への dispatch が責務。

## 判断

### ① 新規モジュール `internal/billing/` を新設、 Webhook handler を移設

```
internal/billing/
  billing.go                          # Module facade
  internal/
    handler/
      webhook_handler.go              # POST /billing/webhook (Stripe)
```

責務: payment + purchase の publicfunctions を import し、 Stripe Webhook event を両 context に dispatch するオーケストレーター。

### ② 命名は業務概念 (`billing`) で、 ベンダー名 (`stripe`) ではない

ADR 0013 ⑤ で「将来 PAY.JP 等の決済 SaaS 追加に備える」 と明言しており、 抽象境界は `WebhookProcessor` interface 側に置く方針。 ここで `internal/stripe/` を切ると:

- PAY.JP を追加した時 `internal/payjp/` ができて、 同じ「申込完了 → Pro 化」 のオーケストレーション責務が複数モジュールに散る
- bounded context (業務分割) ではなく vendor (技術分割) で線を引いたことになる

→ モジュール名は業務 (`billing`)、 vendor 名は handler ファイル名にのみ残す:

```
billing/internal/handler/stripe_webhook_handler.go   # 将来 PAY.JP 追加時:
billing/internal/handler/payjp_webhook_handler.go    # ファイル単位で並ぶ
```

MVP では 1 vendor (Stripe) のみなので、 ファイル名は `webhook_handler.go` で十分。

### ③ publicfunctions の拡張方針

| module | 新規追加 publicfunctions |
|---|---|
| payment | `PaymentWebhookCommand` interface (`CompletePayment` / `CancelPayment` / `RenewPayment` / `HandleFailure`) + `StripeWebhookProcessor` interface (`ParseAndVerify`) |
| purchase | `PurchaseCommand` interface (`ActivateSubscription` 等) |

- 既存 `PaymentCommand` (InitiatePayment 単独) はそのまま、 webhook 系は別 interface に分離 (関心の分離)
- billing module は `paymentCommand` を import せず、 `paymentWebhookCommand` と `stripeWebhookProcessor` だけ受け取る
- 「**外から呼ばれる usecase だけ publicfunctions に昇格**」 のルールは維持。 結果として payment はすべての usecase が外部公開されるが、 集約・repository・状態遷移ロジックは internal 配下のまま隠蔽される (カプセル化は破れていない)

### ④ Composition Root での組み立て順

```go
paymentModule  := payment.NewModule(...)            // 何も知らない
purchaseModule := purchase.NewModule(paymentModule.Command(), ...)   // payment を import
billingModule  := billing.NewModule(
    paymentModule.WebhookCommand(),
    paymentModule.StripeProcessor(),
    purchaseModule.Command(),
)

mux.Handle("/billing/", billingModule.Handler)
```

循環依存は解消。 依存グラフ: `billing → {payment, purchase}`、 `purchase → payment`、 `payment → ()`。

### ⑤ ルーティング変更

| 旧 | 新 |
|---|---|
| `POST /payment/webhook` | `POST /billing/webhook` |

Stripe Dashboard 側の Webhook URL も切り替える必要がある (本番デプロイ時)。 開発時は `.env` の Stripe Webhook 設定を更新。

### ⑥ TX 境界 (再確認)

ADR 0014 ① と ADR 0018 ② の間に「単一 TX で全部 atomic」 vs 「各 usecase が独自 TX」 の解釈差があった。 billing 移設に合わせて以下を確定する:

- **billing handler は usecase を順次呼ぶだけ**、 自前で TX を張らない
- **各 usecase が独自 TX を張る** (ADR 0018 ② 通り)
- 失敗時のシナリオ:
  - `CompletePayment` 成功 → `ActivateSubscription` 失敗の場合、 Stripe リトライで再受信
  - `CompletePayment` は `stripe_events` で冪等 (UNIQUE 違反 = no-op)
  - `ActivateSubscription` は `FindByPaymentID` 既存チェックで冪等 (ADR 0014 ③)
  - 2 回目の Webhook で `CompletePayment` は no-op、 `ActivateSubscription` だけ実行される → 整合
- 「単一 TX で全部 atomic」 は将来の最適化として残す (ADR 0014 ① の図はそちらを想定)。 MVP では「各 usecase 独自 TX + 冪等性で結果整合」 で許容

## なぜ Webhook handler を payment に置き続けないか

代替案: 現状のまま `payment/internal/handler/` に置き、 ADR 0014 末尾の「handler 例外」 で正当化する。

却下理由:
- payment ディレクトリから purchase publicfunctions を import するのは、 grep / 静的解析で「payment が purchase に依存している」 と見える。 ADR の表明と実装が乖離する
- 組み立て順の循環は handler 例外では解消しない
- 将来「Stripe 以外の billing 系 webhook」 (例: 領収書送信、 fraud notification) が増えた時に payment の責務膨張を招く

## なぜ `billing` 内に orchestrator usecase を作らないか

代替案: `billing/internal/usecase/process_stripe_event.go` を作り、 handler は parse + dispatch、 usecase が payment + purchase を順次呼ぶ。

却下理由 (MVP では):
- billing handler 自身が ADR 0014 で言及された「両 context を呼ぶオーケストレーター」 そのもの
- usecase を挟むと、 各 event_type ごとに薄い wrapper を重ねるだけになり構造の冗長化
- handler の責務は「HTTP → アプリ内 command への dispatch」 で SRP 上問題ない

ただし、 将来「冪等性検証ロジック (idempotency_keys) を usecase に閉じたい」 等の理由が出たら usecase 化を再検討。

## やらないこと

- 単一 TX で payment + purchase + outbox を atomic に保存する仕組み (ADR 0014 ① の図を完全実装)
- billing 内に orchestrator usecase 層を作る (上記理由)
- Stripe Customer Portal の return route の置き場所決定 (別 ADR で扱う)

## 影響

### 既存ファイル
- `server/internal/payment/payment.go`: `Handler` field 削除、 `WebhookCommand()` / `StripeProcessor()` getter 追加、 `Command()` (InitiatePayment) は据え置き
- `server/internal/payment/interface/publicfunctions/command.go`: `PaymentWebhookCommand` interface 追加 (別ファイルでも可)
- `server/internal/payment/interface/publicfunctions/`: `StripeWebhookProcessor` interface 追加 (新ファイル `processor.go`)
- `server/internal/payment/internal/handler/webhook_handler.go`: **billing に移設して削除** (人間が移設・実装)
- `server/internal/purchase/interface/publicfunctions/command.go`: `PurchaseCommand` interface 新設 (`ActivateSubscription` 等)
- `server/cmd/server/main.go`: `mux.Handle("/payment/", ...)` 削除、 `mux.Handle("/billing/", ...)` 追加、 billing wire 追加

### 新規ファイル
- `server/internal/billing/billing.go` (Module facade)
- `server/internal/billing/internal/handler/webhook_handler.go` (handler 本体)
- `server/internal/purchase/internal/usecase/activate_subscription.go` (人間が実装)

### ADR 更新
- ADR 0013 末尾に「Webhook handler の置き場所は ADR 0019 で billing に変更」 と追記
- ADR 0014 末尾に同上
- ADR 0018 末尾に「handler は billing 配下に移設」 と追記 (usecase 分割方針は据え置き)

## 関連 ADR

- [ADR 0013](0013-purchase-payment-separation.md): purchase / payment 分離
- [ADR 0014](0014-webhook-idempotency-and-retry.md): Webhook 冪等性
- [ADR 0017](0017-stripe-integration-details.md): Stripe 統合の細部
- [ADR 0018](0018-webhook-usecase-orchestration.md): Webhook usecase 分割
