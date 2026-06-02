// Package traininginfra は Training 集約(Training → TrainingExercise → TrainingSet)を
// MySQL に永続化する。 集約ルートに対する単一 Repository、 子・孫専用は持たない。
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

const upsertTrainingSQL = `
INSERT INTO trainings (id, user_id, started_at, ended_at, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    started_at = VALUES(started_at),
    ended_at   = VALUES(ended_at),
    memo       = VALUES(memo),
    updated_at = VALUES(updated_at)
`

const upsertTrainingExerciseSQL = `
INSERT INTO training_exercises (id, training_id, name, display_order, rest_seconds, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

const upsertTrainingSetSQL = `
INSERT INTO training_sets (id, training_exercise_id, set_number, weight_kg, reps, rest_seconds, memo, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// Save は集約全体を永続化する。
//
// 戦略: 親 upsert → 子孫を全削除 → domain の最新スナップショットを再 INSERT。
// 子の並び替え / 削除 / 追加を「Save 1 回」 で安全に表現するため、 永続化は常にスナップショット差し替えで扱う。
// (呼び出し側が dbtx.Processing で包めば全 Exec が同一 TX で走る。)
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

	// CASCADE で training_sets も連鎖削除される
	if _, err := q.Exec("DELETE FROM training_exercises WHERE training_id = ?", idBytes); err != nil {
		return err
	}

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

// FindByID は集約全体を取り出す。 親 + 子一括 + 孫一括の 3 クエリ。
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
	// 「見つからない」 は UseCase で nil 判定する設計
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	exercises, err := r.loadExercises(ctx, [][]byte{idBytes})
	if err != nil {
		return nil, err
	}
	return toTraining(row, exercises[id.String()])
}

// FindAllByUserIDWithOffsetPagination は一覧版。 親ページング + 子一括 + 孫一括 + COUNT の計 4 クエリ。
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

	var rows []TrainingModel
	_, err = q.Select(&rows,
		"SELECT id, user_id, started_at, ended_at, memo, created_at, updated_at FROM trainings WHERE user_id = ? ORDER BY started_at DESC LIMIT ? OFFSET ?",
		userIDBytes, int32(limit), int32(offset),
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	total, err := q.SelectInt(
		"SELECT COUNT(*) FROM trainings WHERE user_id = ?", userIDBytes,
	)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

	paginator := pagination.NewOffsetPaginator(int(total), offset, limit)

	if len(rows) == 0 {
		return []*trainingdomain.Training{}, paginator, nil
	}

	trainingIDBytes := make([][]byte, 0, len(rows))
	for _, row := range rows {
		trainingIDBytes = append(trainingIDBytes, row.ID)
	}
	exercises, err := r.loadExercises(ctx, trainingIDBytes)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}

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

// loadExercises は親 ID 群に紐づく子(exercises)と孫(sets)を IN 句で一括取得する。
// 「兄弟集合をまとめて取って map で stitch」 が N+1 回避の核。
// 戻り値は trainingID(文字列)→ []*TrainingExercise の map。
func (r *trainingRepository) loadExercises(
	ctx context.Context,
	trainingIDs [][]byte,
) (map[string][]*trainingdomain.TrainingExercise, error) {
	if len(trainingIDs) == 0 {
		return map[string][]*trainingdomain.TrainingExercise{}, nil
	}
	q := dbtx.Querier(ctx, r.dbmap)

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

	exerciseIDBytes := make([][]byte, 0, len(rows))
	for _, row := range rows {
		exerciseIDBytes = append(exerciseIDBytes, row.ID)
	}
	setsByExercise, err := r.loadSets(ctx, exerciseIDBytes)
	if err != nil {
		return nil, err
	}

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

func toTraining(row TrainingModel, exercises []*trainingdomain.TrainingExercise) (*trainingdomain.Training, error) {
	trainingID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.TrainingID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	memo, err := sqlconv.NewString1000FromNullString(row.Memo)
	if err != nil {
		return nil, err
	}
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
	exerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.ID)
	if err != nil {
		return nil, err
	}
	trainingID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.TrainingID](row.TrainingID)
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
	setID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.SetID](row.ID)
	if err != nil {
		return nil, err
	}
	exerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.TrainingExerciseID)
	if err != nil {
		return nil, err
	}
	setNumber, err := valueobject.NewNonNegativeInt(int(row.SetNumber))
	if err != nil {
		return nil, err
	}
	weightKg, err := valueobject.NewNonNegativeDecimalFromString(row.WeightKg)
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
