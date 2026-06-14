# ADR 0012: プレミアム機能の方針 (Pro プラン、 価格、 採用 SaaS)

## ステータス
採用 (2026-06-09)

## コンテキスト

musclead は SODA 入社 (決済基盤チーム配属) の準備を兼ねた個人学習プロジェクト。 本格的な決済 + サブスクリプション機構を実装することで:
- SODA で問われる Webhook / 冪等性 / 監査ログ / 状態遷移 等の経験を積む
- 個人プロダクトとして「課金できる SaaS」 を完成させる

README で「no premium upsell」 を謳っていたが、 SODA 入社準備の文脈で方針転換し、 課金機能を本格的に組み込む。

## 判断

### ① Pro 機能は 1 つに絞る: AI 写真認識

複数機能を並行実装すると焦点がボケる + 学習時間が分散する。 「**決済の深掘り + 非同期処理**」 を SODA 学習の本筋と位置づけ、 Pro 機能はそれを正当化する最小スコープに絞る。

| 候補 | 採用 / 不採用 | 理由 |
|---|---|---|
| **AI 写真認識** (食事カロリー下書き) | ⭕ 採用 | 非同期 / 外部 API / コスト管理 を扱える、 ユーザー価値も高い |
| 昨日の食事コピー | ✗ Free に下ろす候補 | 軽すぎる、 Pro に置く意味薄い |
| バーコードスキャン | ✗ v2 以降 | フロント実装が重く決済本筋から逸れる |

将来 v2 で 2/3 番目を追加する余地は残す (ADR 別途で起票)。

### ② 価格: 月額 480 円 (税込)、 月額のみ、 無料トライアルなし

| 項目 | 値 | 理由 |
|---|---|---|
| 月額 | **480 円 (税込)** | ワンコイン感、 個人 SaaS の現実的価格帯 |
| 表記 | 税込み | 景品表示法の総額表示義務、 Stripe `unit_amount = 480` |
| 年額 | なし (MVP) | 2 つの Price 管理 + プラン切替 (proration) を避ける、 将来追加可能 |
| 無料トライアル | なし (MVP) | trialing / past_due 等の state machine を MVP から外す、 v1.1 以降検討 |

### ③ 決済 SaaS: Stripe

| 観点 | 評価 |
|---|---|
| ドキュメントの厚さ | Stripe が国内最強クラス |
| Test Mode の使いやすさ | テストカード + Stripe CLI |
| Subscription / Webhook 機能の成熟度 | 業界標準 |
| Customer Portal (解約 UI) | Hosted で実装ゼロ |
| Idempotency-Key の標準採用 | 業界の冪等性パターンを学べる |

将来 PAY.JP 等を追加する余地は[ADR 0013](0013-purchase-payment-separation.md) の Facade 設計で確保。

### ④ Pro 機能の認可: 専用 Pro gate middleware

全 Pro API に共通の `ProGate` middleware を 1 つ作り、 ルーティング登録時に被せる。 詳細は [ADR 0014](0014-webhook-idempotency-and-retry.md) で議論。

```go
mux.Handle("POST /ai/recognize-meal", authMw.Wrap(proGate.Wrap(aiHandler.RecognizeMeal)))
```

### ⑤ README の文言修正

旧: `no premium upsell`
新: `Free forever for the core. Premium for the AI-heavy stuff.`

## やらないこと (本 ADR スコープ外)

- 年額プラン (v2 で再検討)
- 無料トライアル (v1.1 で再検討)
- 法人プラン / 複数プラン
- クーポン / 紹介コード
- 領収書発行

## 影響

### ドキュメント
- README 文言修正
- ランディングページの料金セクション追加

### 既存実装
- 既存の 4 集約 (user / weight / meal / training) には影響なし。 Free 機能としてそのまま継続。

### 新規実装の入口
- [ADR 0013](0013-purchase-payment-separation.md): 購入 / 決済 bounded context 分離
- [ADR 0014](0014-webhook-idempotency-and-retry.md): Webhook 冪等性
- [ADR 0015](0015-outbox-pattern-and-async-mail.md): Outbox + メール非同期
- [ADR 0016](0016-ai-photo-recognition.md): AI 写真認識
- [ADR 0017](0017-stripe-integration-details.md): Stripe 統合の細部

## 関連 ADR

- [ADR 0013](0013-purchase-payment-separation.md): 購入と決済の bounded context 分離
