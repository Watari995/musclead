# ADR 0016: AI 写真認識 (食事カロリー下書き生成)

## ステータス
採用 (2026-06-09)

## コンテキスト

[ADR 0012](0012-premium-features-overview.md) で確定した Pro 機能 = AI 写真認識の詳細設計。

要件:
- ユーザーが食事写真をアップロード → AI が料理名 / カロリー / PFC を推定し下書きとして返却
- ユーザーは下書きを確認・編集して meal 記録に反映
- Pro 限定機能 (Free user はアクセス不可)
- 個人開発の規模では月 100 枚程度の使用想定

## 判断

### ① 採用 Vision API: Claude 3.5 Sonnet Vision

| API | 単価 | 食事認識精度 | 日本語 | レイテンシ |
|---|---|---|---|---|
| **Claude 3.5 Sonnet Vision** ✅ | 約 $0.005/枚 | 高 | ◎ | 2-4 秒 |
| OpenAI GPT-4o | 約 $0.005/枚 | 高 | ◎ | 2-3 秒 |
| Gemini 2.0 Flash | 約 $0.001/枚 | 中 | ◎ | 1-2 秒 |
| AWS Rekognition | $0.001/枚 | 低 (汎用ラベル) | △ | < 1 秒 |

選定理由:
- musclead は Claude Code で開発中 (Anthropic 一貫性)
- 日本料理の認識精度が実測で他社と同等以上
- 構造化出力 (JSON schema 指定) が確実
- SODA 入社後の AI/ML 系業務に直結する経験値

### ② 処理方式: **同期処理** (非同期にしない)

ユーザーは認識結果を待ってから meal 記録に反映するため、 「結果を待つ」 行為が本質的に必要。 非同期化しても UX 改善はない一方、 実装複雑度だけ上がる。

| 処理方式 | Pro 化遅延 | 実装複雑度 |
|---|---|---|
| **同期** ✅ | 2-4 秒 (Vision API 待機分のみ) | 低 |
| 非同期 (SQS + Lambda + SSE) | 2-4 秒 (実質同じ) | 高 |

ALB タイムアウト (60 秒) 以内に収まるため、 同期で問題なし。

### ③ 画像アップロード: S3 Presigned URL

```
1. クライアント → POST /ai/upload-url → Presigned URL 取得
2. クライアント → PUT (S3 直接) → S3
3. クライアント → POST /ai/recognize-meal-photo (image_path) → 認識結果
```

musclead は既に profile 画像で Presigned URL パターン採用済み ([ADR 0009](0009-profile-image-storage.md))。 流儀を踏襲。

### ④ DB スキーマ: `ai_recognition_logs` (汎用、 履歴用)

```sql
CREATE TABLE ai_recognition_logs (
  id               BINARY(16)   NOT NULL,
  user_id          BINARY(16)   NOT NULL,
  recognition_type VARCHAR(50)  NOT NULL,        -- 'meal_photo' / 'workout_form' (将来) / 'body_composition' (将来)
  image_path       VARCHAR(500) NOT NULL,        -- S3 path、 musclead 流儀
  result           JSON         NULL,             -- 成功時の認識結果
  error_message    TEXT         NULL,             -- 失敗時
  created_at       DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at       DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (id),
  KEY idx_user_created (user_id, created_at),
  KEY idx_type (recognition_type),
  CONSTRAINT fk_ai_logs_user FOREIGN KEY (user_id) REFERENCES users(id)
);
```

汎用化の理由:
- 将来「ワークアウト判定」「体型分析」 等を同じ Vision API で実装可能
- 共通の処理パイプライン (S3 → Vision API → DB) を再利用できる

(対して `emails` は SES / Twilio / FCM と異なる SaaS を使うので分離。 同じ「汎用化」 でも判断基準が異なる)

### ⑤ Pro 機能ゲート

[ADR 0014](0014-webhook-idempotency-and-retry.md) で確立した `ProGate` middleware を適用:

```go
mux.Handle("POST /ai/recognize-meal-photo",
    authMw.Wrap(proGate.Wrap(rateLimit.Wrap(aiHandler.RecognizeMeal))))
```

middleware の実行順:
1. auth (ユーザー認証)
2. Pro gate (`subscriptions.expires_at > NOW()` チェック)
3. rate limit (本 ADR、 後述)
4. handler

### ⑥ Rate limit: 月 100 枚 (rolling 30 days)

Pro ユーザーといえど無制限にすると Vision API コストが暴走するリスクがあるため、 月 100 枚の上限を設ける。

```sql
SELECT COUNT(*) FROM ai_recognition_logs
WHERE user_id = ?
  AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
  AND error_message IS NULL  -- 失敗はカウントしない
```

| 項目 | 値 |
|---|---|
| 期間 | rolling 30 days |
| 上限 | 100 枚 |
| 計算方式 | DB COUNT クエリ (専用 counter テーブル不要) |
| 失敗時のレスポンス | 429 Too Many Requests |
| 性能 | `idx_user_created` で効率的、 月 100 程度の SELECT は誤差 |

Redis (`ElastiCache`) は [ADR 0010](0010-elasticache-and-fail-open-cache.md) で destroy 済みのため使わず、 DB 集計で十分。

### ⑦ エラーハンドリング

| シナリオ | 対応 |
|---|---|
| Vision API 一時障害 | usecase 内で 1 回リトライ、 失敗時は ai_recognition_logs に error_message を残してユーザーに「認識失敗、 手動入力してください」 |
| Vision API レスポンス JSON 解析失敗 | 同上 (recognition_logs に raw レスポンスを残す) |
| 画像が S3 にない | 400 (クライアントの問題) |
| 画像サイズ超過 (10MB 超) | 413 Payload Too Large (S3 Presigned URL の制約で事前防止) |

## なぜ Claude 3.5 Sonnet Vision なのか

代替案:
- **OpenAI GPT-4o**: 同等性能だが、 musclead は Claude エコシステムで開発。 統一性を取る
- **Gemini 2.0 Flash**: 最安だが精度ばらつきあり、 日本料理データセットでの精度評価が不安定
- **AWS Rekognition**: 食事認識に最適化されておらず、 汎用ラベルしか返さない (「food」 等の粗い分類)
- **自前モデル (SageMaker)**: 個人開発で運用負担が大きすぎる

## なぜ非同期化しないか (再掲)

[ADR 0015](0015-outbox-pattern-and-async-mail.md) でメール送信は非同期化しているのに、 AI 認識を同期にする理由:

| 処理 | ユーザーが待つか | 非同期化の利点 |
|---|---|---|
| メール送信 | 待たない (送れたか気にしない) | ⭕ 失敗時リトライ、 SES と疎結合 |
| AI 認識 | 待つ (結果を見て meal 記録に反映) | ✗ どっちにせよ待つので非同期にする本質的価値なし |

「待つ機能」 を非同期化するのはアンチパターン (ユーザー体感を改善しないのに複雑度だけ上がる)。

## なぜ汎用テーブルなのか (emails は特化、 ai_recognition は汎用)

| 観点 | emails (特化) | ai_recognition (汎用) |
|---|---|---|
| 将来の拡張 | 別 SaaS (SMS = Twilio, push = FCM) | 同じ Vision API で対象が違うだけ |
| 処理パイプライン | チャネルごとに別 | 共通 (S3 → Vision API → DB) |
| 結果データ構造 | 全く違う | type で result JSON が変わる程度 |

→ パイプラインを共有できるかどうかで判断。

## やらないこと

- 結果のキャッシング (同じ画像を 2 度認識しない最適化): YAGNI
- バッチ処理 (複数枚同時アップロード): v2
- 認識精度の継続的評価 / モデル切替: 当面 Claude 任せ
- ワークアウト判定 / 体型分析: v2 で recognition_type を追加
- 月 100 枚超のリクエスト追加課金: 当面シンプルに 429 返す
- レスポンスのストリーミング (token-by-token): 同期で待つので不要

## 影響

### 新規 module
- `internal/ai_recognition/` (Claude Vision client、 認識 usecase、 履歴)

### 新規テーブル
- `ai_recognition_logs`

### 新規ファイル
- `ai_recognition/internal/domain/claude_vision_client.go` (interface)
- `ai_recognition/internal/infra/claude_vision_client.go` (implementation)
- `ai_recognition/internal/usecase/recognize_meal_photo.go`
- `ai_recognition/internal/usecase/check_quota.go` (rate limit)
- `ai_recognition/internal/handler/recognition_handler.go`
- `shared/middleware/pro_gate.go`

### 環境変数 / SSM Parameter
- `ANTHROPIC_API_KEY` (本番は SSM Parameter Store)

### コスト試算
- Pro user 1 人 × 月 100 枚: $0.50
- Pro user 10 人 × 月 100 枚: $5
- → Pro 月額 480 円 (`~$3.2`) で採算ライン 1 人 = 月 9 枚程度

## 関連 ADR

- [ADR 0009](0009-profile-image-storage.md): S3 + Presigned URL パターン
- [ADR 0010](0010-elasticache-and-fail-open-cache.md): ElastiCache 不採用の経緯
- [ADR 0012](0012-premium-features-overview.md): プレミアム機能の方針
