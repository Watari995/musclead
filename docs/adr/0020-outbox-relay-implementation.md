# ADR 0020: outbox-relay の実装方針 (ゼロ円 / ECS relay / DynamoDB 冪等)

## ステータス

採用 (2026-06-15)。 [ADR 0015](0015-outbox-pattern-and-async-mail.md) の transport / 冪等の詳細を **supersede** (Outbox パターンの核心は 0015 を踏襲)。

## コンテキスト

[ADR 0015](0015-outbox-pattern-and-async-mail.md) で「Outbox + SNS + SQS + Lambda でメール送信を非同期化」 と決めた。 実装にあたり、 次の2点を反映して具体化・一部見直す:

1. **追加費用を確実にゼロにする** (個人開発。 常時起動リソース = NAT Gateway / VPC エンドポイント / ElastiCache を一切作らない)
2. **通知の冪等を SODA 流の「短命ストア + TTL」** で実現する (通知ログを永続テーブルに溜めない)

前提: `outbox_events` テーブル / `publishedAt` / `FindPending` / `MarkPublished` / TX 内 INSERT は **0015 で実装済み** (CompletePayment 等が business 更新と同 TX で outbox を書く)。 本 ADR は「**台帳を配る relay と consumer**」 の方針。

## 判断

### ① relay は ECS server 内の goroutine (polling 一本化)

0015 の「Webhook handler から即時 publish + 別 relay の二段構え」 をやめ、 **relay 一本** にする。

- ECS server 起動時に worker goroutine を立て、 一定間隔 (例 10s) で `FindPending(limit)` → 送信 → `MarkPublished` → Save
- relay を ECS 内に置く理由: 既存 ECS は **Stripe API 用の外向き出口を既に持つ** → SQS 送信に NAT 追加不要 = **追加費用ゼロ**。 新規常時起動 compute も作らない

### ② SNS を省略し relay → SQS 直結

消費者は当面メール1種のみ。 fan-out が不要なので **SNS を挟まず SQS に直接送る**。 将来購読者が増えたら SNS を前段に追加する (無料枠内)。

### ③ Lambda は VPC の外 + メッセージ自己完結

- **Lambda を VPC 内に入れない** (入れると RDS 到達のために NAT/VPC エンドポイントが必要 = 有料)。 これが「ゼロ円」 の肝
- そのため Lambda は **DB に触れない**。 メール送信に必要な情報 (宛先 email / イベント種別) は **relay が enrich して SQS メッセージ payload に詰める** (内部 outbox → 外部メッセージの ACL 変換は relay の責務)
- Lambda が叩く先 (SQS / DynamoDB / SES) は全て **VPC 不要の公開 API**

### ④ 通知の冪等は DynamoDB + TTL の条件付き書き込み

非同期は二重実行が避けられない (後述)。 メール送信は冪等でないため明示的に重複を防ぐ。

- **Lambda は送信前に DynamoDB へ event_id を条件付き書き込み** (`attribute_not_exists` = put-if-not-exists)
  - 書けた (初回) → SES 送信
  - 弾かれた (既存) → スキップ
- レコードは **TTL (例 1h) で自動削除** → 通知ログを溜めない (SODA の「揮発 DB + 短い TTL」 と同思想)
- **なぜ DynamoDB か**: Lambda はステートレス + 同時実行のため、 重複は別インスタンスが並列処理しうる。 「送った印」 は **全インスタンスが共有でき・原子的に勝者を1人に決められる外部ストア** が必須。 ElastiCache は VPC 内 = NAT 有料、 DynamoDB は **VPC 外から直接叩けて無料枠内** = ゼロ円を維持できる唯一の選択
- 条件付き書き込みは atomic なので「read してから write」 の競合を起こさない

### ⑤ 重複が起きる3経路 (at-least-once 前提)

④で吸収する重複の発生源:

1. **relay の再送**: 「SQS 送信 → published_at セット」 の間で crash → 次ポーリングで再送
2. **SQS の at-least-once**: 標準キューは稀に同一メッセージを2回配信
3. **Lambda の再試行**: SES 送信成功後、 SQS へ完了を返す前に crash → 再配信

→ いずれも「副作用」 と「完了記録」 が原子的にできないことに起因。 ④の DynamoDB 冪等で「ちょうど1回」 に収束させる。

### ⑥ SES はドメイン検証 (DKIM)

- 送信元 `no-reply@musclead.com`。 **ドメイン検証** (Route53 に DKIM CNAME を追加、 Terraform で作成)
- 理由: アドレス検証だと1アドレスからしか送れず DKIM 無しで迷惑メール判定されやすい。 ドメイン検証は任意のアドレス + DKIM 署名で到達率が高く、 Route53 管理なので追加作業も小さい
- 注意: 新規 SES は **サンドボックス** (検証済み宛先のみ送信可)。 一般ユーザー宛は **production access 申請** が必要

### ⑦ 初回対象イベント

**payment_succeeded (申込み完了メール) のみ**。 更新 (renewed) / 解約 (canceled) は後続フェーズ。

## ゼロ円の根拠

| リソース | 無料枠 (毎月・恒久) | 本件の実使用 |
|---|---|---|
| SQS | 100万 req | 月数件 → $0 |
| Lambda | 100万実行 + 40万 GB秒 | 月数件 → $0 |
| DynamoDB | 25 WCU/RCU + 25GB | 月数件 → $0、 TTL 削除も無料 |
| SES | 後述 | 月数通 → ほぼ $0 ($0.10/1000通) |

**常時起動リソース (NAT / VPC エンドポイント / ElastiCache) を作らない** ことが条件。 全て pay-per-use + always-free 枠。

## 影響

- 新規 Terraform: SQS + DLQ / Lambda (VPC 外) / DynamoDB (TTL 有効) / SES domain identity + DKIM / IAM
- ECS server に relay worker goroutine + SQS 送信 adapter + enrich (user email 取得) を追加
- Lambda は別デプロイ単位 (Go)。 SQS トリガー
- 0015 の「即時 publish」「emails 永続テーブルで冪等」 は本 ADR で置換 (relay polling / DynamoDB+TTL)

## 更新履歴

- 2026-06-15: 新規作成。 0015 の transport (SNS 省略・relay を ECS goroutine 一本化) と冪等 (DynamoDB+TTL) を具体化・supersede。
