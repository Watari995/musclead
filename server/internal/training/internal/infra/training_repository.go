// Package traininginfra は Training 集約(Training → TrainingExercise → TrainingSet)を
// MySQL に永続化する。 1集約 = 1 Repository の原則に従い、 子・孫専用の Repository は持たない。
//
// 読み書きは dbtx.Querier 経由なので、 呼び出し側が TransactionManager.Processing で
// 包めば集約全体が自動的に1 TX に収まる。
package traininginfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlquery"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type trainingRepository struct {
	dbmap *gorp.DbMap
}

func NewTrainingRepository(dbmap *gorp.DbMap) trainingdomain.TrainingRepository {
	return &trainingRepository{dbmap: dbmap}
}

// 親 = trainings は upsert(更新可能フィールドだけ ON DUPLICATE KEY UPDATE)。
const upsertTrainingSQL = `
INSERT INTO trainings (id, user_id, started_at, ended_at, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    started_at = VALUES(started_at),
    ended_at   = VALUES(ended_at),
    memo       = VALUES(memo),
    updated_at = VALUES(updated_at)
`

// 子 = training_exercises は Save 時に「親に紐づく行を全削除 → 再 INSERT」 の戦略を取るので
// INSERT のみで十分。 個別 UPDATE は不要。
const upsertTrainingExerciseSQL = `
INSERT INTO training_exercises (id, training_id, name, display_order, rest_seconds, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

// 孫 = training_sets は親 (exercise) が CASCADE で削除されるので、 親リセット時に
// 一緒に消える。 こちらも INSERT のみ。
const upsertTrainingSetSQL = `
INSERT INTO training_sets (id, training_exercise_id, set_number, weight_kg, reps, rest_seconds, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// Save は集約全体を一括で永続化する。
//
// 戦略:
//  1. 親 trainings を upsert(更新可能フィールドのみ)
//  2. その training_id に紐づく training_exercises を全削除
//     (CASCADE で training_sets も連鎖削除される)
//  3. domain 側に残っている exercises / sets を全件 INSERT し直す
//
// 「子孫を一度全消ししてから再投入」 する理由は、 「子の削除 / 並び替え / 追加」 を
// 1 つの Save で安全に表現するため。 部分更新は集約内部メソッドで完結させ、
// 永続化レイヤは「集約のスナップショットを丸ごと書く」 ことに専念する。
//
// 呼び出し側が dbtx.TransactionManager.Processing で包んでいれば、 ここの全 Exec は
// 同一トランザクションで実行される(dbtx.Querier が ctx の tx を拾うため)。
func (r *trainingRepository) Save(ctx context.Context, training *trainingdomain.Training) error {
	idBytes, err := training.ID().Bytes()
	if err != nil {
		return err
	}
	userIDBytes, err := training.UserID().Bytes()
	if err != nil {
		return err
	}

	q := dbtx.Querier(ctx, r.dbmap)

	// 1. trainings を upsert
	if _, err := q.Exec(upsertTrainingSQL,
		idBytes,
		userIDBytes,
		training.StartedAt(),
		sqlconv.ToNullTime(training.EndedAt()),
		sqlconv.String1000ToNullString(training.Memo()),
		training.CreatedAt(),
		training.UpdatedAt(),
	); err != nil {
		return err
	}

	// 2. 子(exercises)を全削除 → CASCADE で孫(sets)も削除される
	if _, err := q.Exec("DELETE FROM training_exercises WHERE training_id = ?", idBytes); err != nil {
		return err
	}

	// 3. domain 上の最新状態を順に INSERT
	for _, ex := range training.Exercises() {
		exIDBytes, err := ex.ID().Bytes()
		if err != nil {
			return err
		}
		if _, err := q.Exec(upsertTrainingExerciseSQL,
			exIDBytes,
			idBytes,
			ex.Name().Value(),
			int32(ex.DisplayOrder().Value()),
			sqlconv.NonNegativeIntToNullInt32(ex.RestSeconds()),
			sqlconv.String1000ToNullString(ex.Memo()),
			ex.CreatedAt(),
			ex.UpdatedAt(),
		); err != nil {
			return err
		}
		for _, set := range ex.Sets() {
			setIDBytes, err := set.ID().Bytes()
			if err != nil {
				return err
			}
			if _, err := q.Exec(upsertTrainingSetSQL,
				setIDBytes,
				exIDBytes,
				int32(set.SetNumber().Value()),
				set.WeightKg().String(),
				int32(set.Reps().Value()),
				sqlconv.NonNegativeIntToNullInt32(set.RestSeconds()),
				sqlconv.String1000ToNullString(set.Memo()),
				set.CreatedAt(),
				set.UpdatedAt(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

// FindByID は集約全体を取り出す。
// 親 1 件 + 子の一括取得 + 孫の一括取得、 計 3 クエリで完結させる(N+1 を避ける)。
func (r *trainingRepository) FindByID(ctx context.Context, id valueobject.TrainingID) (*trainingdomain.Training, error) {
	idBytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}

	q := dbtx.Querier(ctx, r.dbmap)

	var row TrainingModel
	err = q.SelectOne(&row,
		"SELECT id, user_id, started_at, ended_at, memo, created_at, updated_at FROM trainings WHERE id = ?",
		idBytes,
	)
	// 「見つからない」 は UseCase で nil 判定する設計なので、 ここでは nil, nil を返す。
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 子・孫を取得して map で stitched する(loadExercises 内で sets も解決済み)
	exercises, err := r.loadExercises(ctx, [][]byte{idBytes})
	if err != nil {
		return nil, err
	}
	return toTraining(row, exercises[id.String()])
}

// FindAllByUserIDWithOffsetPagination は一覧版。
// 親一覧 + 子一括 + 孫一括の 3 クエリ + COUNT クエリ、 計 4 クエリ。
//
// ページング後の親 ID 集合をベースに「IN 句で子・孫を1ショットで取る」 ことで
// N+1 にならない作りにしている。
func (r *trainingRepository) FindAllByUserIDWithOffsetPagination(
	ctx context.Context,
	userID valueobject.UserID,
	limit int,
	offset int,
) ([]*trainingdomain.Training, pagination.OffsetPaginator, error) {
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	q := dbtx.Querier(ctx, r.dbmap)

	// 1. 親(trainings)をページングして取得
	var rows []TrainingModel
	_, err = q.Select(&rows,
		"SELECT id, user_id, started_at, ended_at, memo, created_at, updated_at FROM trainings WHERE user_id = ? ORDER BY started_at DESC LIMIT ? OFFSET ?",
		userIDBytes, int32(limit), int32(offset),
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	// 2. 総件数(COUNT)を取って Paginator を作る
	total, err := q.SelectInt(
		"SELECT COUNT(*) FROM trainings WHERE user_id = ?", userIDBytes,
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	paginator := pagination.OffsetPaginator{
		CurrentPage:  offset/limit + 1,
		ItemsPerPage: limit,
		TotalItems:   int(total),
		TotalPages:   (int(total) + limit - 1) / limit,
	}

	// 親 0 件なら子・孫を引きにいく必要なし
	if len(rows) == 0 {
		return []*trainingdomain.Training{}, paginator, nil
	}

	// 3. 親 ID 一覧を作って、 子(exercises)とその孫(sets)を一括ロード
	trainingIDBytes := make([][]byte, 0, len(rows))
	for _, row := range rows {
		trainingIDBytes = append(trainingIDBytes, row.ID)
	}
	exercises, err := r.loadExercises(ctx, trainingIDBytes)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	// 4. 親 row 順に entity を組み立てる
	trainings := make([]*trainingdomain.Training, 0, len(rows))
	for _, row := range rows {
		idStr, err := sqlconv.UUIDStringFromBytes(row.ID)
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		t, err := toTraining(row, exercises[idStr])
		if err != nil {
			return nil, pagination.OffsetPaginator{}, err
		}
		trainings = append(trainings, t)
	}
	return trainings, paginator, nil
}

// DeleteByID は親 trainings を消す。 子・孫は FK の ON DELETE CASCADE で連鎖削除される。
func (r *trainingRepository) DeleteByID(ctx context.Context, id valueobject.TrainingID) error {
	idBytes, err := id.Bytes()
	if err != nil {
		return err
	}
	q := dbtx.Querier(ctx, r.dbmap)
	_, err = q.Exec("DELETE FROM trainings WHERE id = ?", idBytes)
	return err
}

// loadExercises は指定された複数 training_id に紐づく training_exercises を一括取得し、
// さらにその exercises の sets も loadSets で一括取得する。
//
// 戻り値の map は「training_id (string)」 → 「その training に属する exercises」。
// 呼び出し側は親を回しながら map[parentID] でアクセスして組み立てる。
//
// この「IN 句で兄弟集合をまとめて取る → map で束ねる」 パターンが N+1 回避の正体。
func (r *trainingRepository) loadExercises(
	ctx context.Context,
	trainingIDs [][]byte,
) (map[string][]*trainingdomain.TrainingExercise, error) {
	if len(trainingIDs) == 0 {
		return map[string][]*trainingdomain.TrainingExercise{}, nil
	}
	q := dbtx.Querier(ctx, r.dbmap)

	// IN 句用の "?,?,?" と args を共通ヘルパで作る
	placeholders, args := sqlquery.InPlaceholders(trainingIDs)
	var rows []TrainingExerciseModel
	_, err := q.Select(&rows,
		"SELECT id, training_id, name, display_order, rest_seconds, memo, created_at, updated_at FROM training_exercises WHERE training_id IN ("+placeholders+") ORDER BY display_order ASC",
		args...,
	)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return map[string][]*trainingdomain.TrainingExercise{}, nil
	}

	// 孫(sets)もこのタイミングで全部まとめて読む
	exerciseIDBytes := make([][]byte, 0, len(rows))
	for _, row := range rows {
		exerciseIDBytes = append(exerciseIDBytes, row.ID)
	}
	setsByExercise, err := r.loadSets(ctx, exerciseIDBytes)
	if err != nil {
		return nil, err
	}

	// stitched 結果: trainingID(文字列)→ []*TrainingExercise
	result := make(map[string][]*trainingdomain.TrainingExercise, len(rows))
	for _, row := range rows {
		trainingIDStr, err := sqlconv.UUIDStringFromBytes(row.TrainingID)
		if err != nil {
			return nil, err
		}
		exerciseIDStr, err := sqlconv.UUIDStringFromBytes(row.ID)
		if err != nil {
			return nil, err
		}
		ex, err := toTrainingExercise(row, setsByExercise[exerciseIDStr])
		if err != nil {
			return nil, err
		}
		result[trainingIDStr] = append(result[trainingIDStr], ex)
	}
	return result, nil
}

// loadSets は指定 exercise 群に紐づく training_sets を一括取得し、 exercise_id ベースで map に束ねる。
func (r *trainingRepository) loadSets(
	ctx context.Context,
	exerciseIDs [][]byte,
) (map[string][]*trainingdomain.TrainingSet, error) {
	if len(exerciseIDs) == 0 {
		return map[string][]*trainingdomain.TrainingSet{}, nil
	}
	q := dbtx.Querier(ctx, r.dbmap)

	placeholders, args := sqlquery.InPlaceholders(exerciseIDs)
	var rows []TrainingSetModel
	_, err := q.Select(&rows,
		"SELECT id, training_exercise_id, set_number, weight_kg, reps, rest_seconds, memo, created_at, updated_at FROM training_sets WHERE training_exercise_id IN ("+placeholders+") ORDER BY set_number ASC",
		args...,
	)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]*trainingdomain.TrainingSet, len(rows))
	for _, row := range rows {
		exerciseIDStr, err := sqlconv.UUIDStringFromBytes(row.TrainingExerciseID)
		if err != nil {
			return nil, err
		}
		s, err := toTrainingSet(row)
		if err != nil {
			return nil, err
		}
		result[exerciseIDStr] = append(result[exerciseIDStr], s)
	}
	return result, nil
}

// toTraining は DB row + 子 exercises から domain 上の Training entity を組み立てる。
// 各 VO への変換時に検証ロジックが走るので、 異常データは error として上にバブルさせる。
func toTraining(row TrainingModel, exercises []*trainingdomain.TrainingExercise) (*trainingdomain.Training, error) {
	idStr, err := sqlconv.UUIDStringFromBytes(row.ID)
	if err != nil {
		return nil, err
	}
	trainingID, err := valueobject.NewPrimaryIDFromString[valueobject.TrainingID](idStr)
	if err != nil {
		return nil, err
	}
	userIDStr, err := sqlconv.UUIDStringFromBytes(row.UserID)
	if err != nil {
		return nil, err
	}
	userID, err := valueobject.NewPrimaryIDFromString[valueobject.UserID](userIDStr)
	if err != nil {
		return nil, err
	}
	memo, err := sqlconv.NewString1000FromNullString(row.Memo)
	if err != nil {
		return nil, err
	}
	// nil スライス保持を防ぐ。 domain は「空でも初期化済み」 の方が安全。
	if exercises == nil {
		exercises = []*trainingdomain.TrainingExercise{}
	}
	return trainingdomain.NewTraining(
		*trainingID,
		*userID,
		row.StartedAt,
		sqlconv.FromNullTime(row.EndedAt),
		memo,
		row.CreatedAt,
		row.UpdatedAt,
		exercises,
	), nil
}

func toTrainingExercise(row TrainingExerciseModel, sets []*trainingdomain.TrainingSet) (*trainingdomain.TrainingExercise, error) {
	idStr, err := sqlconv.UUIDStringFromBytes(row.ID)
	if err != nil {
		return nil, err
	}
	exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](idStr)
	if err != nil {
		return nil, err
	}
	trainingIDStr, err := sqlconv.UUIDStringFromBytes(row.TrainingID)
	if err != nil {
		return nil, err
	}
	trainingID, err := valueobject.NewPrimaryIDFromString[valueobject.TrainingID](trainingIDStr)
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewString50(row.Name)
	if err != nil {
		return nil, err
	}
	displayOrder, err := valueobject.NewNonNegativeInt(int(row.DisplayOrder))
	if err != nil {
		return nil, err
	}
	restSeconds, err := sqlconv.NewNonNegativeIntFromNullInt32(row.RestSeconds)
	if err != nil {
		return nil, err
	}
	memo, err := sqlconv.NewString1000FromNullString(row.Memo)
	if err != nil {
		return nil, err
	}
	if sets == nil {
		sets = []*trainingdomain.TrainingSet{}
	}
	return trainingdomain.NewTrainingExercise(
		*exerciseID,
		*trainingID,
		*name,
		*displayOrder,
		restSeconds,
		memo,
		row.CreatedAt,
		row.UpdatedAt,
		sets,
	), nil
}

func toTrainingSet(row TrainingSetModel) (*trainingdomain.TrainingSet, error) {
	idStr, err := sqlconv.UUIDStringFromBytes(row.ID)
	if err != nil {
		return nil, err
	}
	setID, err := valueobject.NewPrimaryIDFromString[valueobject.SetID](idStr)
	if err != nil {
		return nil, err
	}
	exerciseIDStr, err := sqlconv.UUIDStringFromBytes(row.TrainingExerciseID)
	if err != nil {
		return nil, err
	}
	exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](exerciseIDStr)
	if err != nil {
		return nil, err
	}
	setNumber, err := valueobject.NewNonNegativeInt(int(row.SetNumber))
	if err != nil {
		return nil, err
	}
	weightKg, err := sqlconv.NewNonNegativeDecimalFromString(row.WeightKg)
	if err != nil {
		return nil, err
	}
	reps, err := valueobject.NewNonNegativeInt(int(row.Reps))
	if err != nil {
		return nil, err
	}
	restSeconds, err := sqlconv.NewNonNegativeIntFromNullInt32(row.RestSeconds)
	if err != nil {
		return nil, err
	}
	memo, err := sqlconv.NewString1000FromNullString(row.Memo)
	if err != nil {
		return nil, err
	}
	return trainingdomain.NewTrainingSet(
		*setID,
		*exerciseID,
		*setNumber,
		*weightKg,
		*reps,
		restSeconds,
		memo,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}
