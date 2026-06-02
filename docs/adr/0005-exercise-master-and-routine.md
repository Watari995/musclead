# ADR 0005: 種目マスタ導入と Routine モジュール配置

## ステータス
採用 (2026-06-02)

## コンテキスト

`training_exercises.name` に種目名を直書きしていたため、 次の問題が顕在化した。

1. **表記揺れ**: 「ベンチプレス」 「ベンチ・プレス」 「Bench Press」 が別物として扱われる
2. **rename 不能**: 過去の training 全件を JOIN なしで一括更新する手段がない
3. **SSOT 違反**: 同じ種目が複数の training_exercises 行に複製される
4. **Routine 未実装**: 「テンプレから training を起こす」 機能を入れたいが、 現状の name 直書きだと routine_items も name 直書きになり、 同じ問題が増殖する

## 決定

### 1. Exercise マスタテーブルを導入する
- ユーザーごとに種目マスタ `exercises(id, user_id, name, UNIQUE(user_id, name))` を持つ
- `training_exercises.name` を `training_exercises.exercise_id`(FK)に差し替え
- routine 系も `routine_items.exercise_id` で参照

### 2. Exercise は独立した集約
- DDD 集約境界: Exercise は自分のライフサイクル(rename / 将来 archive)を持つ
- Training 集約 / Routine 集約からは **exercise_id で参照のみ**(直接埋め込まない)
- TrainingExercise 側は `display_order` / `rest_seconds` / `memo` / `sets` をそのまま保持(セッション固有のためマスタに上げない)

### 3. Routine は Training モジュール内に配置
- Routine 集約 = Routine(ルート)+ RoutineItem(子、 exercise_id のみ持つ)
- training/exercise/routine を **すべて training モジュール内**に同居:
  - `internal/training/internal/domain/{training,training_exercise,training_set,exercise,routine}.go`
  - cross-module 公開 interface 不要、 凝集度が高い

### 4. 書き込み = 集約厳密 / 読み出し = JOIN(CQRS-lite)
- 書き込み: Training/Exercise/Routine それぞれ Repository を分離、 集約境界を守る
- 読み出し: TrainingRepository / RoutineRepository が view 用に exercises を JOIN して name を取得
  - N+1 を回避しつつ、 ドメイン操作の集約境界は守る

### 5. 既存型の rename(衝突回避)
- 既存 `Exercise` → `TrainingExercise`
- 既存 `Set` → `TrainingSet`
- 新規マスタ = `Exercise`
- DB テーブル名(`training_exercises` / `training_sets` / `exercises`)と一致して読みやすい

### 6. Exercise 削除はマスタ側でブロック(当面)
- FK は `ON DELETE RESTRICT`
- 使用中の Exercise は削除不可、 UI でエラー
- 将来 `deleted_at` で論理削除(archive)機能を足す

## 理由

### a. 既存実装を破壊しすぎない
- TrainingExercise の構造(`display_order` / `rest_seconds` / `memo` / `sets`)はそのまま流用
- 差し替えは `name` フィールド → `exercise_id` だけ
- 既存の `training` 集約のロジック・SQL の骨格は維持

### b. マスタテーブルの name で SSOT を確立する
- 種目名の正本は exercises テーブル1箇所のみ
- rename は exercises を1行更新するだけで全 training に反映
- 表記揺れは UNIQUE(user_id, name) で抑制

### c. 過去の値を ID 経由で参照しやすい
- 「前回のベンチプレス、 何 kg だった?」 の検索が `WHERE te.exercise_id = ?` で確実
- 名前検索だと日本語の表記揺れ / 過去 rename によるヒット漏れリスクあり
- インデックスも数値 FK の方が安定して効く

### d. 集約境界の明確化(DDD 原則)
- Exercise の編集と Training の編集が独立する
- 将来 `video_url` / `target_muscles` / `last_used_at` 等を Exercise に追加しても、 Training の集約に変更が及ばない

### e. Routine と Training の整合
- Routine を Training モジュール内に置くことで、 cross-module の公開 interface 不要
- 「routine から training を雛形展開する」 ロジックが同一モジュール内で完結
- Exercise も他モジュール(meal 等)からは参照しないため、 同居が YAGNI 的に適切

### f. 書き込み集約厳密 / 読み出し JOIN
- 書き込み側で集約境界を守ることで整合性・テスト性が保たれる
- 読み出し側を JOIN にすることで N+1 を排除し、 個人開発レベルでもパフォーマンス問題を作らない

## 不採用案

### 案 A: exercise_name 直書きのまま(現状維持)
- 表記揺れ・rename 不能・SSOT 違反のため却下

### 案 B: name で参照(マスタ持つが ID ではなく name で繋ぐ)
- 日本語の表記揺れリスク
- rename したら過去の検索がヒットしなくなる
- 数値 FK の効率に劣る
- 却下

### 案 C: Exercise を独立モジュール化
- 現時点で meal や他モジュールから参照する予定がない
- 公開 interface の追加コストに見合わない(YAGNI)
- 必要になった時点で `interface/publicfunctions/` 経由で切り出す
- 当面は training モジュール内に同居

### 案 D: Routine を独立モジュール化
- Routine は Training と密結合(雛形→記録の関係)
- cross-module call が増えるだけで価値が薄い
- 同居が妥当

### 案 E: Exercise 削除を CASCADE
- 過去の training 履歴が消える → データ消失リスク大、 却下
- RESTRICT + 将来 archive(soft delete)で対応

## 結果

### メリット
- ✅ 種目名の SSOT 確立、 表記揺れ排除
- ✅ rename が1操作で全箇所反映
- ✅ Routine 実装の足場が綺麗(exercise_id で参照するだけ)
- ✅ 集約境界が明確、 将来の Exercise 拡張(video_url 等)も Training に影響しない

### デメリット / 受け入れるトレードオフ
- ❌ 既存 training_exercises データの破棄が必要(dev 環境のため許容、 本番想定なら backfill SQL を別途用意)
- ❌ UX シフト: training 編集の「種目名 text 入力」 → 「Exercise マスタ select」
- ❌ Exercise 削除の制約(当面 FK RESTRICT、 archive 機能は future task)
- ❌ Exercise UseCase / Handler / DTO / 管理 FE 画面の追加実装コスト

## 実施計画

| Step | 内容 | 担当 |
|---|---|---|
| 1 | 既存 `Exercise` → `TrainingExercise`、 `Set` → `TrainingSet` 全箇所 rename(挙動変更なし) | AI |
| 2 | Exercise マスタ実装: migration / domain / infra / usecase / handler / dto | 人間(BE)+ AI(FE 管理画面) |
| 3 | TrainingExercise の `name` → `exercise_id` 移行: migration / domain / infra / dto / FE select 化 | 人間(BE)+ AI(FE) |
| 4 | Routine 実装: migration / domain / infra / usecase / handler / dto / FE | 人間(BE)+ AI(FE) |
| 5 | 前回値 API(`GET /trainings/exercises/{exercise_id}/latest`)+ Routine 使用時の前回値表示 | 人間(BE)+ AI(FE) |

各 Step ごとに commit + push、 Step 間で破壊的変更を分離する。

## 関連 ADR

- [0002: DDD + Modular Monolith](./0002-ddd-modular-monolith.md) — モジュール境界の方針
- [0003: Database conventions](./0003-database-conventions.md) — テーブル命名 / FK / UNIQUE の規約
- [0004: Module public interface](./0004-module-public-interface.md) — モジュール間連携(将来 Exercise を切り出す際に適用)
