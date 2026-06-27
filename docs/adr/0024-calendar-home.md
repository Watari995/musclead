# ADR-0024: カレンダーホーム機能の設計

- ステータス: Accepted
- 日付: 2026-06-27
- 関連: [ADR-0002 DDD + Modular Monolith](0002-ddd-modular-monolith.md)

## コンテキスト

アプリのホーム画面として、トレーニング・食事・体重の記録をカレンダー形式で一覧できる機能を追加する。
記録した日付に点を表示し、日付をタップするとその日のサマリーを確認できる。

また、学習目的として goroutine（`sync.WaitGroup` による並列クエリ）を自然な形で実装する題材でもある。

## 決定

### モバイル画面構成

- ホームタブを5タブ目として一番左に追加し、アプリ起動時の初期タブとする
- 月単位カレンダー表示（ボタン + スワイプで月移動。未来月への移動も可）
- 記録のある日付に色付き点を表示（トレーニング・食事・体重で色分け、複数並列）
- 日付タップ → カレンダー下部にサマリーエリアが展開
- 記録がない日をタップした場合はサマリーエリアを非表示
- 起動時の初期選択日は「今日」
- サマリーの各項目タップ → 既存の詳細画面に遷移

### 点の色設定

ユーザーごとに色をカスタマイズできるよう `user_preferences` テーブルに色カラムを追加する。
色はサーバーから取得し、クライアントはAPIの値をそのまま適用する（ハードコードしない）。

```sql
ALTER TABLE user_preferences
  ADD COLUMN training_color VARCHAR(7) NOT NULL DEFAULT '#4A90E2',
  ADD COLUMN meal_color     VARCHAR(7) NOT NULL DEFAULT '#7ED321',
  ADD COLUMN weight_color   VARCHAR(7) NOT NULL DEFAULT '#FF6B6B';
```

### モジュール配置

cross-domain な集計ロジックを `calendar` モジュールとして新設する。
`home` という命名は UI の概念であるため採用しない。

```
server/internal/calendar/
  calendar.go
  internal/
    handler/
      calendar_handler.go
    usecase/
      get_monthly_summary.go   -- 月の記録フラグ取得（goroutine）
      get_daily_summary.go     -- 日別サマリー取得（goroutine）
    dto/
```

`calendar` モジュールは training・weight・meal モジュールの public interface にのみ依存し、各モジュールの内部には触れない。

### API 設計

| メソッド | パス | 説明 |
|---|---|---|
| `GET` | `/calendar/monthly-summary?year=2026&month=6` | 月の記録フラグ一覧 |
| `GET` | `/calendar/daily-summary?date=2026-06-27` | 日別サマリー詳細 |
| `GET` | `/user/preferences` | 色設定取得 |
| `PUT` | `/user/preferences` | 色設定更新 |

#### `GET /calendar/monthly-summary` レスポンス

```json
{
  "days": [
    {
      "date": "2026-06-01",
      "hasTraining": true,
      "hasMeal": false,
      "hasWeight": false
    }
  ]
}
```

#### `GET /calendar/daily-summary` レスポンス

```json
{
  "trainings": [
    { "trainingId": "xxx", "exerciseCount": 4, "totalSets": 12 }
  ],
  "meals": [
    { "mealId": "xxx", "name": "朝食", "calories": 700, "eatenAt": "07:30" }
  ],
  "weights": [
    {
      "weightId": "xxx",
      "weightKg": 72.3,
      "bodyFatPercentage": 18.5,
      "skeletalMuscleKg": 32.1,
      "measuredAt": "07:00"
    }
  ]
}
```

### goroutine 実装方針

`monthly-summary` / `daily-summary` ともに training・weight・meal の3クエリを並列実行する。

**パターン: `sync.WaitGroup` + 個別エラーハンドリング**

`errgroup` を使わない理由: 1つのクエリが失敗しても他の成功データを返すため（部分失敗許容）。
各 goroutine は自身の結果とエラーを独立した変数に書き込み、`WaitGroup.Wait()` 後にまとめてレスポンスを組み立てる。

### タイムゾーン

日付の境界は JST（UTC+9）基準とする。クエリの日付フィルタは JST → UTC 変換を行ってから DB に渡す。

### パフォーマンス目標

`monthly-summary` のレスポンスタイムは 0.5 秒以内を目標とする。

## 代替案

- **`home` モジュールとして新設**: フロントエンドの概念をバックエンドに持ち込むため不採用。
- **`errgroup` で並列処理**: 1つ失敗で全キャンセルになり部分失敗許容ができないため不採用。
- **色設定をクライアント側に持つ**: アプリ削除で設定が消えるため不採用。

## 影響

- migration: `user_preferences` に色カラム3つを追加
- `server/internal/calendar/` を新規作成し `main.go` に追加
- `server/internal/user/` に preferences の GET/PUT エンドポイントを追加
- Mobile: ホームタブを追加、カレンダーウィジェット実装
