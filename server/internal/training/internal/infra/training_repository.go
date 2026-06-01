package traininginfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
	"github.com/shopspring/decimal"
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
		toNullString(training.Memo()),
		training.CreatedAt(),
		training.UpdatedAt(),
	); err != nil {
		return err
	}

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
			toNullInt32FromNonNegativeInt(ex.RestSeconds()),
			toNullString(ex.Memo()),
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
				toNullInt32FromNonNegativeInt(set.RestSeconds()),
				toNullString(set.Memo()),
				set.CreatedAt(),
				set.UpdatedAt(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

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

	paginator := pagination.OffsetPaginator{
		CurrentPage:  offset/limit + 1,
		ItemsPerPage: limit,
		TotalItems:   int(total),
		TotalPages:   (int(total) + limit - 1) / limit,
	}

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

func (r *trainingRepository) DeleteByID(ctx context.Context, id valueobject.TrainingID) error {
	idBytes, err := id.Bytes()
	if err != nil {
		return err
	}
	q := dbtx.Querier(ctx, r.dbmap)
	_, err = q.Exec("DELETE FROM trainings WHERE id = ?", idBytes)
	return err
}

func (r *trainingRepository) loadExercises(
	ctx context.Context,
	trainingIDs [][]byte,
) (map[string][]*trainingdomain.TrainingExercise, error) {
	if len(trainingIDs) == 0 {
		return map[string][]*trainingdomain.TrainingExercise{}, nil
	}
	q := dbtx.Querier(ctx, r.dbmap)

	placeholders, args := makePlaceholders(trainingIDs)
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

	placeholders, args := makePlaceholders(exerciseIDs)
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

func makePlaceholders(items [][]byte) (string, []interface{}) {
	placeholders := ""
	args := make([]interface{}, 0, len(items))
	for i, item := range items {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, item)
	}
	return placeholders, args
}

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
	memo, err := nullStringToString1000(row.Memo)
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
	restSeconds, err := nullInt32ToNonNegativeInt(row.RestSeconds)
	if err != nil {
		return nil, err
	}
	memo, err := nullStringToString1000(row.Memo)
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
	weightDecimal, err := decimal.NewFromString(row.WeightKg)
	if err != nil {
		return nil, err
	}
	weightKg, err := valueobject.NewNonNegativeDecimal(weightDecimal)
	if err != nil {
		return nil, err
	}
	reps, err := valueobject.NewNonNegativeInt(int(row.Reps))
	if err != nil {
		return nil, err
	}
	restSeconds, err := nullInt32ToNonNegativeInt(row.RestSeconds)
	if err != nil {
		return nil, err
	}
	memo, err := nullStringToString1000(row.Memo)
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

func toNullString(v *valueobject.String1000) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: v.Value(), Valid: true}
}

func toNullInt32FromNonNegativeInt(v *valueobject.NonNegativeInt) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: int32(v.Value()), Valid: true}
}

func nullStringToString1000(v sql.NullString) (*valueobject.String1000, error) {
	if !v.Valid {
		return nil, nil
	}
	s, err := valueobject.NewString1000(v.String)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func nullInt32ToNonNegativeInt(v sql.NullInt32) (*valueobject.NonNegativeInt, error) {
	if !v.Valid {
		return nil, nil
	}
	n, err := valueobject.NewNonNegativeInt(int(v.Int32))
	if err != nil {
		return nil, err
	}
	return n, nil
}
