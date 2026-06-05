# ADR 0006: Routine 集約の構造と Training 展開フロー

## ステータス
採用 (2026-06-03)

## コンテキスト

ADR 0005 で Routine を training モジュール内の独立集約として配置する方針を決めたが、
具体的な構造と「Training に展開する」 ユースケースを未確定にしていた。 Phase 3 として
実装に着手するにあたり、 集約の入れ子・展開フロー・編集の伝播範囲を明文化する。

主な検討点:

1. **どの粒度まで Routine が抱えるか**: 種目順だけか、 各種目のデフォルトセット(重量・レップ)まで持つか
2. **Training への展開時に何を copy するか**: 編集後の Routine が既存 Training に影響するか
3. **expand を Backend で集約横断ユースケースとして実装するか、 FE が orchestrate するか**

## 決定

### 1. 集約は 2 階層 (Routine → RoutineExercise)

```
Routine (集約ルート)
  ├── id, user_id, name
  └── exercises: []*RoutineExercise
       ├── id
       ├── routine_id
       ├── exercise_id  (Exercise マスタへの FK)
       └── display_order
```

`RoutineExercise` は **`exercise_id` と `display_order` のみ**を持つ。
- `rest_seconds` / `memo` / sets 等は持たない
- 「セットや重量はトレーニング当日に決める」 という musclead の運用モデルに沿う
- 集約構造が浅くなり、 Save 戦略が単純化する

### 2. Training への展開は「種目順をコピーするだけ」

「ルーティン A から training を作成」 アクションは以下を行う:

1. Routine の `exercises` を `display_order` 順で取得
2. 各 `RoutineExercise.exercise_id` を **新しい TrainingExercise** にコピー(`sets = []` 空)
3. `Training.started_at = now`、 `ended_at = null`、 `memo = null` で新規作成
4. 作成後、 ユーザーは Training 編集画面でセットを埋める

### 3. 展開後は完全独立(Copy on use)

- Training の中の TrainingExercise は **Routine と参照関係を持たない**(`routine_id` を保存しない)
- Routine 編集 / 削除は既存 Training に影響しない
- 「過去のセッションが時系列で正本」 という設計思想を維持

### 4. expand は FE が orchestrate する

「ルーティンから training を作成」 専用エンドポイントは Backend に作らない。 FE が:

1. `GET /routines/{id}` で routine と exercise 順を取得
2. `POST /trainings` でその exercise_id 列を `exercises[].exercise_id` に詰めて投げる
3. レスポンスの `training_id` で `/trainings/{id}/edit` に遷移

集約横断の orchestration を usecase 層に持ち込まず、 各集約は自前のユースケースだけ提供する。

### 5. Exercise マスタへの参照は ON DELETE RESTRICT

`routine_exercises.exercise_id` も `training_exercises.exercise_id` と同じく
`ON DELETE RESTRICT`。 routine から参照されている Exercise を削除すると
**409 Conflict + `exercise_used_in_training_error`**(*) が返る。

(*) FK 違反のエラー判定が「training から参照されてる」 と「routine から参照されてる」
の両方を 1451 で捕まえるため、 既存の `ExerciseUsedInTrainingError` を共用する。
将来的に「どこから参照されているか」 を見せたくなったら、 error code を一般化して
`ExerciseInUseError` に rename する余地を残す。

### 6. Routine 自身の所有者チェックは FindByIDAndUserID パターン

Exercise と同じく `RoutineRepository.FindByIDAndUserID` で所有者検証を data access 層
に集約する。 usecase は NotFound チェック 1 回で済み、 Permission 分岐を持たない。

## 理由

a. **「セットは当日決める」 派 vs 「テンプレ重量も含めて再現する」 派**: musclead は
   筋肥大トラッキングを主目的とし、 重量レップは前回比較で決めるため、 ルーティンに
   重量情報を持つ価値は薄い。 種目順だけテンプレ化できれば 8 割の要件を満たす。

b. **集約境界の整合性**: RoutineExercise が `exercise_id + display_order` だけなら、
   Routine 集約は値オブジェクトに近い軽量さで Save 戦略が単純(親 upsert → 子全削除 →
   再 INSERT)。

c. **expand を FE 主導にする**: 「ルーティンを使う = 種目順だけコピー」 という単純さは
   FE で十分扱える。 Backend に集約横断ユースケースを増やすほどの抽象化は不要。

d. **Copy on use**: 過去 Training を「その日に決めた構成のスナップショット」 として
   保つことで、 履歴比較や履歴可視化が壊れない。

e. **FK RESTRICT 流用**: Routine からの参照も RESTRICT にすることで、 user が「使用中
   の Exercise を間違って削除」 する事故を物理的に防ぐ。

## 不採用案

A. **RoutineExercise にデフォルトセット (target_weight, target_reps, sets) を持たせる**
   - 棄却: 重量はその日に決める運用と合わない。 musclead は「ルーティン = 種目順テンプレ」
     とシンプルに割り切る。

B. **Training が `routine_id` を保持して参照関係を残す**
   - 棄却: Routine 編集が過去 Training に「いつの間にか反映される」 事故が起きる。
     履歴の不変性を優先。

C. **「ルーティンから training を作成」 専用エンドポイント (POST /routines/{id}/start) を作る**
   - 棄却: Backend が集約横断 usecase(Routine 読 + Training 書)を持つことになり、 集約境界が
     曖昧になる。 FE で 2 リクエストするのと、 学習価値・ロジックの集中度のトレードオフで、
     現状は FE 主導が筋。 将来 endpoint 化したくなれば 1 ユースケース足すだけ。

D. **Routine を集約 1 階層 (`name + exercise_ids []string`) で JSON 列に持つ**
   - 棄却: 個別 RoutineExercise を id で扱えなくなり、 並び替え UI / アナリティクス /
     将来の拡張で詰む。 リレーショナルに展開しておく。

## 結果

- migration `000010_create_routines.up.sql` + `000011_create_routine_exercises.up.sql` を追加
- domain: `routine.go`, `routine_exercise.go`, `routine_repository.go` を新設
- infra: `routine_models.go`, `routine_repository.go` を新設
- usecase: `create_routine.go`, `update_routine.go`, `delete_routine.go`, `find_routine.go`,
  `list_routines.go` の 5 本(expand 専用ユースケースは作らない)
- handler: `routine_handler.go` を新設、 training モジュールの Facade に組み込み
- DTO: `routine_dto.go`
- swag annotation + swag init
- FE: `/routines` (一覧)、 `/routines/new` (作成)、 `/routines/[id]/edit` (編集)、
  `/routines/[id]` (詳細 + 「このルーティンで開始」 ボタンが FE 内で POST /trainings を叩く)
- Header に「ルーティン」 nav リンク追加

## 実施計画

1. migration 2 本(routines + routine_exercises、 RESTRICT FK 付き)
2. domain entity + Spec + repository interface
3. infra model + repository(Save は親 upsert → 子全削除 → 再 INSERT のスナップショット差し替え)
4. usecase 5 本
5. handler + dto + swag init
6. training モジュール Facade (`training.go`) で Routine 用 DI 配線、 `Module.RoutineHandler` 追加
7. `main.go` で `/routines` route を mount
8. FE: schema 再生成、 Routine ページ 4 つ + Header nav + 「このルーティンで開始」 アクション
