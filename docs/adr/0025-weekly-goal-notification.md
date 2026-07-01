# ADR-0025: 週次目標通知機能の設計

- ステータス: Accepted
- 日付: 2026-06-29
- 関連: [ADR-0002 DDD + Modular Monolith](0002-ddd-modular-monolith.md)、[ADR-0024 カレンダーホーム](0024-calendar-home.md)

## コンテキスト

ユーザーが週次の目標（トレーニング回数・平均カロリー・体重変化）を設定し、毎週日曜にバックグラウンドワーカーが達成状況を判定してアプリ内通知を生成する機能を追加する。

また、学習目的として Go の goroutine と channel（buffered channel による producer/consumer、worker pool）を自然な形で実装する題材でもある。ADR-0024 で `sync.WaitGroup` と `errgroup` を学んだ延長として、channel を用いた worker pool パターンを習得する。

## 決定

### モジュール構成

`user_weekly_goals` はユーザー設定の一部として `user` モジュールに統合する。`goal` を独立モジュールとするほど大きな関心事ではなく、`user_preferences` と同じ文脈で管理するのが自然なため。

```
server/internal/user/          -- 既存。user_weekly_goals を追加
server/internal/notification/  -- アプリ内通知の生成・既読管理（新規）
server/cmd/goal-checker/       -- 週次チェック worker（goroutine/channel 本体）（新規）
```

### テーブル設計

**`user_weekly_goals`**（ユーザーに1行、upsert）

| カラム | 型 | 説明 |
|---|---|---|
| id | BINARY(16) | PK |
| user_id | BINARY(16) | FK → users、UNIQUE |
| training_count | INT NULL | 週のトレーニング目標回数 |
| calorie_average | INT NULL | 週の平均カロリー目標(kcal) |
| weight_change_kg | DECIMAL(4,1) NULL | 週の体重変化目標(kg、符号付き) |
| created_at | DATETIME(6) | |
| updated_at | DATETIME(6) | |

**`notifications`**

| カラム | 型 | 説明 |
|---|---|---|
| id | BINARY(16) | PK |
| user_id | BINARY(16) | FK → users |
| notification_type | VARCHAR(50) | `'weekly_goal'` など |
| metadata | JSON | 型固有データ（例: `{"goal_type":"training","is_achieved":true,"actual":4,"target":3}`） |
| read_at | DATETIME(6) NULL | NULL = 未読 |
| created_at | DATETIME(6) | |

### API 設計

| メソッド | パス | 説明 |
|---|---|---|
| `GET` | `/user/weekly-goal` | 自分の目標取得 |
| `PUT` | `/user/weekly-goal` | 目標更新（未設定は null） |
| `GET` | `/notifications` | 通知一覧（`unread_count` 含む） |
| `GET` | `/notifications/:id` | 通知詳細 |
| `PUT` | `/notifications/:id/read` | 既読化 |

### 通知文言

| goal_type | 達成 | 未達成 |
|---|---|---|
| training | 今週のトレーニングは{actual}回でした。目標の{target}回を達成しました！ | 今週のトレーニングは{actual}回でした。目標の{target}回に届きませんでした。 |
| calorie | 今週の平均カロリーは{actual}kcalでした。目標の{target}kcal以内を達成しました！ | 今週の平均カロリーは{actual}kcalでした。目標の{target}kcalを超えています。 |
| weight | 今週の体重変化は{actual}kgでした。目標の{target}kgを達成しました！ | 今週の体重変化は{actual}kgでした。目標に届きませんでした。 |

### goroutine / channel 実装方針

週1回（日曜23時）の ticker が発火した時点で全ユーザーを対象にチェックを行う。

```
ticker(日曜23時)
  └─ producer goroutine: 全ユーザーIDを chan UserID（buffered）に流す
       └─ worker goroutine × N（固定数）: channel から受け取り
            └─ 各ユーザーの3目標を errgroup で並列チェック
                 └─ notifications テーブルに書き込む
```

| パターン | 使用箇所 |
|---|---|
| buffered channel | producer → worker 間のジョブキュー |
| worker pool | goroutine 数を N に固定してDB負荷を制御 |
| errgroup | ユーザーごとの3目標並列チェック |
| done channel（context） | graceful shutdown |

### 達成判定

- **training**: `実際の回数 >= target`
- **calorie**: `実際の平均 <= target`
- **weight_change_kg**: `目標と同方向かつ絶対値が target 以上`（例: 目標 -0.5kg → 実際 -0.6kg は達成、-0.3kg は未達成）

目標が null のユーザーはその種別をスキップし通知しない。

### UI フロー（モバイル）

```
ホーム画面ヘッダー右: 🔔 [未読数バッジ]
  └─ タップ → 通知一覧画面（未読は背景ハイライト）
       └─ タップ → 通知詳細画面（達成/未達成アイコン + 文言全文）
            └─ 詳細取得と同時に既読化（PUT /notifications/:id/read を連続呼び出し）
```

## 代替案

- **push 通知（FCM/APNs）**: インフラが未整備のため不採用。アプリ内通知で代替。
- **メール通知（Resend）**: 週1通知でメールは過剰。アプリ内通知の方が UX として自然。
- **毎日チェック + 重複排除ロジック**: 週1判定で十分なため複雑性を避けて不採用。
- **EAV（goal_type / target_value の行持ち）**: 型安全性が下がるため不採用。カラム持ちで明示的に管理する。
- **`goal` を独立モジュールとする**: `user_preferences` と同じくユーザー設定の一部であるため `user` モジュールに統合。独立させるほどの大きさではない。

## 影響

- migration: `user_weekly_goals`、`notifications` テーブルを新規作成
- `server/internal/user/` に `user_weekly_goals` のドメイン・infra・usecase・handler を追加
- `server/internal/notification/` を新規作成し `main.go` にルート追加
- `server/cmd/goal-checker/` に worker を新規作成
- Mobile: ホームタブヘッダーにベルアイコン追加、通知一覧・詳細画面を新規作成
