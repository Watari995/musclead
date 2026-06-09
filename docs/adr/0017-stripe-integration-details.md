# ADR 0017: Stripe 統合の細部 (Checkout、 Customer Portal、 Customer 管理)

## ステータス
採用 (2026-06-09)

## コンテキスト

[ADR 0012](0012-premium-features-overview.md) で Stripe 採用、 [ADR 0013](0013-purchase-payment-separation.md) で抽象化レイヤー、 [ADR 0014](0014-webhook-idempotency-and-retry.md) で Webhook 処理を決定済み。 本 ADR は Stripe 統合の実務的な細部 (Dashboard セットアップ、 UI 構成、 Customer 管理、 開発フロー) を扱う。

## 判断

### ① 決済 UI: Stripe Checkout (Hosted page)

Stripe が提供するホスト型決済ページに `window.location` でリダイレクト。

```
musclead UI 「Pro 申込」 ボタン
  → POST /purchase/subscribe → Stripe Checkout Session 作成 → URL 返却
  → window.location = url
  → Stripe Checkout ページでカード入力 → 決済
  → success_url (= musclead UI) にリダイレクト
  → Webhook で musclead 内に subscription 有効化
```

#### 代替案との比較

| 案 | フロント実装 | PCI DSS | 学習軸 (決済本筋) |
|---|---|---|---|
| **A. Stripe Checkout (Hosted)** ✅ | 最小 | Stripe 完全管理 | Webhook / 冪等性に集中可 |
| B. Stripe Payment Element (Embedded) | 中程度 | SAQ-A | UI 実装に時間取られる |
| C. Stripe Elements (個別 input) | 重い | 中 | musclead 規模で過剰 |

学習目的は「Webhook / 冪等性 / outbox / 非同期」 という決済の**本筋**にあり、 UI 実装に時間を吸われると本末転倒。 Checkout で決済 UI 部分を完全に Stripe に委譲。

### ② サブスク管理 UI: Stripe Customer Portal

ユーザーの**解約・カード変更**は Stripe Customer Portal を利用 (自前実装しない)。

```go
// POST /purchase/portal-session
func (h *PortalHandler) CreateSession(w, r) {
    payment, _ := h.paymentQuery.FindLatestByUserID(ctx, userID)
    session, _ := h.stripeClient.CreatePortalSession(ctx, stripe.PortalSessionParams{
        Customer:  payment.StripeCustomerID,
        ReturnURL: "https://app.musclead.com/settings/subscription",
    })
    writeJSON(w, map[string]string{"portal_url": session.URL})
}
```

#### Portal で公開する操作 (MVP)

| 操作 | 公開する? |
|---|---|
| 解約 (`cancel_at_period_end`) | ⭕ |
| カード変更 | ⭕ |
| 請求書履歴 | ✗ MVP 不要 |
| メールアドレス変更 | ✗ MVP 不要 (musclead UI で行う) |
| プラン変更 | ✗ 1 プランしかない |

### ③ 解約方針: 期末解約 (cancel_at_period_end)

```
ユーザーが Portal で解約ボタン押下 (7/9)
  → Stripe: cancel_at_period_end = true
  → subscriptions.canceled_at = 7/9, status = 'canceled'
  → expires_at = 7/15 (期末まで Pro 機能アクセス可)
  → 7/15 経過後、 Pro gate middleware が expires_at < NOW() で 403 を返す
  → Stripe Webhook 'customer.subscription.deleted' で subscriptions.status = 'expired'
```

#### Pro 判定ロジック

```go
SELECT subscriptions WHERE user_id = ?
  ORDER BY created_at DESC LIMIT 1
// status に関係なく、 expires_at だけで判定
return sub.ExpiresAt.After(time.Now())
```

→ `status = canceled` でも期末までは Pro 扱い (Stripe デフォルト挙動と一致)。

### ④ Stripe Customer 管理: payments.stripe_customer_id を再利用

同一 user が解約 → 再加入する時、 過去の `stripe_customer_id` を再利用して **二重 Customer 作成を防ぐ**。

```go
func (uc *InitiatePayment) Execute(ctx, input) {
    existing, _ := uc.paymentRepo.FindLatestByUserID(ctx, input.UserID)
    var customerID string
    if existing != nil && existing.StripeCustomerID != "" {
        customerID = existing.StripeCustomerID   // ⭐ 既存再利用
    } else {
        customer, _ := stripeClient.CreateCustomer(ctx, ...)
        customerID = customer.ID
    }

    session, _ := stripeClient.CreateCheckoutSession(ctx, stripe.CheckoutSessionParams{
        Customer: stripe.String(customerID),    // ⭐ 既存 Customer に紐付け
        Mode:     stripe.String("subscription"),
    })
}
```

#### ダブルタップ (連打) による二重作成防止

[ADR 0014](0014-webhook-idempotency-and-retry.md) の `idempotency_keys` middleware で:
- 1 回目: `idempotency_keys` INSERT (processing) → Stripe.CreateCustomer 実行
- 2 回目 (ダブルタップ): UNIQUE 制約違反 → 既存レコード SELECT → "処理中" → 409 Conflict
- → Stripe.CreateCustomer が 2 回呼ばれない

→ **payment_customers 専用テーブルは作らない** (`payments.stripe_customer_id` で十分)。

### ⑤ /users/me に subscription 情報を統合

フロントは `/users/me` 1 リクエストで「Pro かどうか」「期限はいつか」 を取得できる。 user handler が `purchase.SubscriptionQuery` を呼んで合成。

```json
GET /users/me
{
  "user": { "id": "...", "email": "..." },
  "subscription": {
    "plan": "pro",
    "status": "active",
    "expires_at": "2026-07-09T00:00:00Z"
  }
}
```

handler レイヤーでの複数 context 統合は依存方向違反ではない (handler は最外層)。

### ⑥ Stripe Dashboard セットアップ (手動、 1 度きり)

| 項目 | 設定方法 |
|---|---|
| Product / Price | Stripe Dashboard で手動作成、 Price ID を SSM Parameter に保存 |
| Webhook Endpoint | musclead 本番 URL を登録、 Webhook signing secret を SSM に保存 |
| Customer Portal | Dashboard の Configuration で「解約」「カード変更」 を有効化 |
| Test Mode → Live Mode 切替 | 本番公開直前に切替 |

Terraform で Stripe リソースを管理する選択肢もあるが、 非公式 Provider のメンテが薄い + 滅多に変えないため手動で十分。

### ⑦ 開発時の Stripe CLI 利用

ローカル開発時:

```bash
# 1. Stripe CLI で Webhook を localhost に転送
stripe listen --forward-to localhost:8080/payment/webhook

# 2. テストカードで決済テスト
# Stripe Checkout: 4242 4242 4242 4242

# 3. Webhook イベントを擬似発火
stripe trigger checkout.session.completed
stripe trigger invoice.payment_failed
stripe trigger customer.subscription.deleted
```

| テストカード | 用途 |
|---|---|
| `4242 4242 4242 4242` | 成功 |
| `4000 0000 0000 9995` | 即時 decline (残高不足) |
| `4000 0000 0000 0341` | サブスク後の月次更新失敗 |
| `4000 0027 6000 3184` | 3D Secure 認証 |

### ⑧ API Key の管理

| Key | 環境変数 / 保存場所 |
|---|---|
| Publishable Key (`pk_test_xxx` / `pk_live_xxx`) | フロント環境変数 (Next.js `NEXT_PUBLIC_STRIPE_PK`) |
| Secret Key (`sk_test_xxx` / `sk_live_xxx`) | サーバ環境変数、 本番は SSM Parameter Store |
| Webhook Signing Secret (`whsec_xxx`) | サーバ環境変数、 本番は SSM Parameter Store |

既存 musclead の SSM 統合パターンを踏襲。

## なぜ Stripe Checkout なのか (Payment Element でなく)

- 学習軸が「決済本筋」 (Webhook / 冪等性 / outbox) にあり、 UI 実装はノイズ
- Stripe Checkout / Payment Element どちらでも Webhook 以降の処理は同じ
- フロント実装が 1 日で済み、 AI 写真認識 ([ADR 0016](0016-ai-photo-recognition.md)) の実装時間を確保できる
- ブランド一貫性は MVP では二の次

将来本格的に musclead を有料公開する時、 Payment Element に移行する余地は残す (Webhook 処理は変えなくて済む)。

## なぜ Customer Portal なのか (自前実装でなく)

- 解約・カード変更 UI を自前で作ると 1-2 週間吸われる
- Stripe Portal は無料、 月次更新で機能拡充される
- musclead のブランドが Stripe ページに飛ぶことを許容するレベルの個人プロダクト
- 業界事例も多数

将来「引き止め機能」「キャンセル理由ヒアリング」 等を実装する時、 自前 UI に切替検討。

## なぜ Customer を payments.stripe_customer_id で管理するか (専用テーブルでなく)

代替案: `payment_customers` テーブル (`user_id` UNIQUE で 1 user 1 customer 強制)。

却下理由:
- `payments.stripe_customer_id` を SELECT して再利用するだけで十分
- 同時並行による二重作成は `idempotency_keys` で吸収可能
- テーブル増加コストに見合う価値が低い (個人開発、 法人プラン等の複雑要件なし)

法人プラン追加時 (= 1 user に複数 subscription を持たせたい時) は専用テーブル化を再検討。

## やらないこと

- 領収書 / 請求書発行 (Stripe Dashboard で確認可能、 ユーザー直アクセスは Portal で)
- クーポン / 紹介コード
- 法人プラン
- 自前カード入力 UI (Payment Element / Elements)
- Stripe Tax (税計算、 月額 0.5% 追加)
- ローカルでの Stripe API モック化 (Test Mode で十分)

## 影響

### 新規エンドポイント
- `POST /purchase/subscribe` (idempotency middleware 適用)
- `POST /purchase/portal-session` (Customer Portal 起動)
- `POST /payment/webhook` (Stripe Webhook 受信)

### Stripe Dashboard セットアップ (本番デプロイ前に手動)
- Product「Pro Plan」 + Price「480 JPY / month」
- Webhook Endpoint `https://api.musclead.com/payment/webhook`
- Customer Portal Configuration (解約 + カード変更)

### SSM Parameter (本番)
- `/musclead/stripe/secret_key`
- `/musclead/stripe/webhook_signing_secret`
- `/musclead/stripe/price_id` (Pro plan の Price ID)

### フロント
- `web/src/app/settings/subscription/page.tsx` (Free → 申込ボタン、 Pro → Customer Portal ボタン)
- `web/src/features/purchase/api/subscribe.ts` (Idempotency-Key 付き fetch)
- 申込成功ページ `/subscription/success` (Pro 化処理中の polling 表示)

### Test Mode 開発期間中
- Stripe Test Mode を使用
- ローカル開発時は Stripe CLI で Webhook 転送

## 関連 ADR

- [ADR 0012](0012-premium-features-overview.md): プレミアム機能の方針
- [ADR 0013](0013-purchase-payment-separation.md): purchase / payment 分離
- [ADR 0014](0014-webhook-idempotency-and-retry.md): Webhook 冪等性
