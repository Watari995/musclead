package traininginfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlerr"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type exerciseRepository struct {
	dbmap *gorp.DbMap
}

func NewExerciseRepository(dbmap *gorp.DbMap) trainingdomain.ExerciseRepository {
	return &exerciseRepository{dbmap: dbmap}
}

const insertExerciseSQL = `
INSERT INTO exercises (id, user_id, name, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
`

const updateExerciseSQL = `
UPDATE exercises
SET name = ?, updated_at = ?
WHERE id = ?
`

func (r *exerciseRepository) FindByIDAndUserID(ctx context.Context, id valueobject.ExerciseID, userID valueobject.UserID) (*trainingdomain.Exercise, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row ExerciseModel
	err = q.SelectOne(&row,
		"SELECT id, user_id, name, created_at, updated_at FROM exercises WHERE id = ? and user_id = ?",
		idBytes, userIDBytes,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toExercise(row)
}

func (r *exerciseRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit int, offset int) ([]*trainingdomain.Exercise, pagination.OffsetPaginator, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := userID.Bytes()
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	var rows []ExerciseModel
	_, err = q.Select(&rows,
		"SELECT id, user_id, name, created_at, updated_at FROM exercises WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		bytes, limit, offset)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	total, err := q.SelectInt("SELECT COUNT(*) FROM exercises WHERE user_id = ?", bytes)
	if err != nil {
		return nil, pagination.OffsetPaginator{}, err
	}
	paginator := pagination.NewOffsetPaginator(int(total), offset, limit)
	if len(rows) == 0 {
		return []*trainingdomain.Exercise{}, paginator, nil
	}
	exercises := make([]*trainingdomain.Exercise, 0, len(rows))
	for _, row := range rows {
		e, err := toExercise(row)
		if err != nil {
			return nil, paginator, err
		}
		exercises = append(exercises, e)
	}
	return exercises, paginator, nil
}

func (r *exerciseRepository) Save(ctx context.Context, exercise *trainingdomain.Exercise) error {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := exercise.ID().Bytes()
	if err != nil {
		return err
	}
	// update first
	result, err := q.Exec(updateExerciseSQL, exercise.Name().Value(), exercise.UpdatedAt(), idBytes)
	if err != nil {
		if sqlerr.IsDuplicateKey(err) {
			return myerror.NewExerciseNameAlreadyExistsError()
		}
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// if not updated, insert
	if rowsAffected == 0 {
		params, err := buildInsertExerciseParams(exercise)
		if err != nil {
			return err
		}
		_, err = q.Exec(insertExerciseSQL, params...)
		if sqlerr.IsDuplicateKey(err) {
			return myerror.NewExerciseNameAlreadyExistsError()
		}
		return err
	}
	return nil
}

func (r *exerciseRepository) DeleteByID(ctx context.Context, id valueobject.ExerciseID) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := id.Bytes()
	if err != nil {
		return err
	}
	if _, err := q.Exec("DELETE FROM exercises WHERE id = ?", bytes); err != nil {
		return err
	}
	return nil
}

func toExercise(row ExerciseModel) (*trainingdomain.Exercise, error) {
	exerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewString50(row.Name)
	if err != nil {
		return nil, err
	}
	return trainingdomain.NewExercise(*exerciseID, *userID, *name, row.CreatedAt, row.UpdatedAt), nil
}

func buildInsertExerciseParams(exercise *trainingdomain.Exercise) ([]any, error) {
	bytes, err := exercise.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := exercise.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	return []any{bytes, userIDBytes, exercise.Name().Value(), exercise.CreatedAt(), exercise.UpdatedAt()}, nil
}
