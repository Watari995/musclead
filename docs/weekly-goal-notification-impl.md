# 週次目標達成通知 — 実装メモ

> PR: `feat/weekly-goal-notification`
> 目的: goroutine/channel の学習を兼ねて、週次目標達成チェックワーカーとアプリ内通知機能を追加する。

---

## やったこと

### DB migration
- `000028_create_user_weekly_goals.up.sql` — ユーザーごとの週次目標テーブル
- `000029_create_notifications.up.sql` — 通知テーブル

### valueobject
- `NotificationType` — 通知種別VO (`weekly_goal` 等)
- `WeightChangeKg` — 体重変化VO（+-値を許容する `DecimalBase` ベース）
- `UserWeeklyGoalID` / `NotificationID` — 各エンティティのPrimaryID VO

### user module（`server/internal/user/`）
- `UserWeeklyGoal` エンティティ追加（週次目標: トレーニング回数・カロリー平均・体重変化）
- `GET /users/me/weekly-goal` / `PUT /users/me/weekly-goal` API
- public interface に `GetWeeklyGoal` / `GetAllUserIDs` を追加（workerが使用）

### notification module（`server/internal/notification/`）新規
- `Notification` エンティティ（type・metadata・read_at）
- `GET /notifications` — 一覧（未読数付き）
- `GET /notifications/{id}` — 詳細
- `PATCH /notifications/{id}/read` — 既読化
- public interface `NotificationCommand.Create` — worker から通知を作成するための窓口

### training / meal / weight module — 週次集計クエリ追加
| module | メソッド | 戻り値 |
|---|---|---|
| training | `CountSessionsByWeek(ctx, userID, weekStart)` | `NonNegativeInt` |
| meal | `GetAverageCaloriesInAWeek(ctx, userID, weekStart)` | `*NonNegativeDecimal` |
| weight | `GetWeightChangeInAWeek(ctx, userID, weekStart)` | `*WeightChangeKg` |

各モジュールで domain interface → infra(SQL) → usecase → publicfunctions の順に実装済み。

### goal-checker worker（`server/cmd/goal-checker/`）骨格のみ
- `main.go` — DB接続・module初期化・`run()` 起動
- `worker.go` — buffered channel + worker pool + ticker scheduler（骨格実装済み）
- `checker.go` — `checkAndNotify` 関数（**TODO: 未実装**）

---

## これから実装すること

### バックエンド

#### 1. `checker.go` の `checkAndNotify` を実装（最重要・学習メイン）

```go
func checkAndNotify(ctx, userID, weekStart, ...) error {
    // 1. 週次目標を取得
    goal, err := userQuery.GetWeeklyGoal(ctx, userpublicfunctions.GetWeeklyGoalInput{UserID: userID})
    // goal が nil（未設定）なら return nil

    // 2. 各種実績を取得
    count   := trainingQuery.CountSessionsByWeek(ctx, userID, weekStart)
    calorie := mealQuery.GetAverageCaloriesInAWeek(ctx, userID, weekStart)
    change  := weightQuery.GetWeightChangeInAWeek(ctx, userID, weekStart)

    // 3. 目標と実績を比較
    //    - goal.TrainingCount != nil && count < goal.TrainingCount → 未達成
    //    - goal.CalorieAverage != nil && calorie != nil → 比較
    //    - goal.WeightChangeKg != nil && change != nil → 比較

    // 4. metadata を組み立てて通知作成
    metadata := valueobject.Metadata{
        "training_goal":   goal.TrainingCount,
        "training_actual": count,
        "achieved":        achieved,
        // ...
    }
    return notifCommand.Create(ctx, userID, notificationType, metadata)
}
```

通知の `notification_type` は `weekly_goal` 固定。
`metadata` の中身でフロントが表示文言を組み立てる想定。

#### 2. worker の graceful shutdown
現在 `ctx.Done()` でworkerを止めているが、OS signal（SIGTERM/SIGINT）を受け取って cancel する処理が必要。

```go
// main.go に追加
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
go func() {
    <-sigCh
    cancel()
}()
```

#### 3. worker の deploy 設定
Oracle VM 上で goal-checker バイナリを systemd サービスとして動かす、または cron で週1起動する設定が必要。
現状 `deploy-vm.yml` は server バイナリのみを対象にしているため、goal-checker のビルド・転送・起動を追加する。

---

### モバイル（Flutter / `mobile/`）— AI担当

#### 実装すべき画面・機能

**1. ホームタブ — ベルアイコン（未読バッジ）**
- AppBar右上にベルアイコンを追加
- `GET /notifications` を呼び出して `unread_count` を取得
- 未読が1件以上あればバッジ表示（赤丸 + 数字）
- タップで通知一覧画面に遷移

**2. 通知一覧画面**
- `GET /notifications` の結果を一覧表示
- 各行: 通知種別アイコン + 日時 + 既読/未読の視覚的区別（未読は太字など）
- タップで通知詳細画面に遷移

**3. 通知詳細画面**
- `GET /notifications/{id}` で詳細取得
- 画面表示と同時に `PATCH /notifications/{id}/read` で既読化
- `metadata` の内容をパースして表示文言を組み立てる
  - `weekly_goal` タイプの場合: 目標値・実績値・達成/未達成を表示

**4. 週次目標設定画面**
- `GET /users/me/weekly-goal` で現在の目標を取得
- `PUT /users/me/weekly-goal` で目標を更新
- フィールド: トレーニング回数（nullable整数）・カロリー平均（nullable整数）・体重変化kg（nullable小数）
- 設定なし（null）は「目標なし」として扱う

#### metadata の表示文言（weekly_goal）
```
// metadata 構造
{
  "training_goal":    3,        // 目標回数（null=未設定）
  "training_actual":  2,        // 実績回数
  "calorie_goal":     2000,     // 目標カロリー平均（null=未設定）
  "calorie_actual":   1850.5,   // 実績カロリー平均（null=データなし）
  "weight_goal":      -1.0,     // 目標体重変化kg（null=未設定）
  "weight_actual":    -0.5,     // 実績体重変化kg（null=データなし）
  "achieved":         false     // 全目標達成フラグ
}
```

表示例:
- 達成: 「今週の目標を達成しました 🎉」
- 未達成: 「今週のトレーニングは 2/3 回でした」

---

## goroutine/channel の学習まとめ

このワーカーで使っているGoの並行処理パターン:

| パターン | 使用箇所 | 概要 |
|---|---|---|
| buffered channel | `ch := make(chan UserID, N)` | 送受信を疎結合にするキュー |
| worker pool | `for i := 0; i < N; i++ { go func() {...} }` | N本のgoroutineがchannelを並列受信 |
| `for range ch` | workerのループ | `close(ch)` されると自動でループ終了 |
| `close(ch)` | ctx.Done()時 | 全workerへの終了合図 |
| `sync.WaitGroup` | `wg.Wait()` | 全worker終了を待つ |
| ticker | `time.NewTicker(7 * 24 * time.Hour)` | 週1回のスケジューリング |
| select | tickerとctx.Doneの分岐 | 複数channelの同時監視 |
