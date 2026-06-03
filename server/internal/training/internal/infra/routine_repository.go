package traininginfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/shared/sqlerr"
	"github.com/Watari995/musclead/internal/shared/sqlquery"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type routineRepository struct {
	dbmap *gorp.DbMap
}

func NewRoutineRepository(dbmap *gorp.DbMap) trainingdomain.RoutineRepository {
	return &routineRepository{dbmap: dbmap}
}

const insertRoutineSQL = `
INSERT INTO routines (id, user_id, name, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
`

const insertRoutineExerciseSQL = `
INSERT INTO routine_exercises (id, routine_id, exercise_id, display_order, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
`

const updateRoutineSQL = `
UPDATE routines
SET name = ?, updated_at = ?
WHERE id = ?
`

func (r *routineRepository) FindByIDAndUserID(ctx context.Context, id valueobject.RoutineID, userID valueobject.UserID) (*trainingdomain.Routine, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row RoutineModel
	err = q.SelectOne(&row,
		"SELECT id, user_id, name, created_at, updated_at FROM routines WHERE id = ? AND user_id = ?",
		idBytes, userIDBytes,
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

	return toRoutine(row, exercises[id.String()])
}

func (r *routineRepository) Save(ctx context.Context, routine *trainingdomain.Routine) error {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := routine.ID().Bytes()
	if err != nil {
		return err
	}
	userIDBytes, err := routine.UserID().Bytes()
	if err != nil {
		return err
	}
	// update first なければ insertにする
	result, err := q.Exec(updateRoutineSQL, routine.Name().Value(), routine.UpdatedAt(), idBytes)
	if err != nil {
		if sqlerr.IsDuplicateKey(err) {
			return myerror.NewRoutineNameAlreadyExistsError()
		}
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// if not updated, insert
	if rowsAffected == 0 {
		result, err = q.Exec(insertRoutineSQL, idBytes, userIDBytes, routine.Name().Value(), routine.CreatedAt(), routine.UpdatedAt())
		if err != nil {
			if sqlerr.IsDuplicateKey(err) {
				return myerror.NewRoutineNameAlreadyExistsError()
			}
			return err
		}
	}

	// 子を削除
	if _, err := q.Exec("DELETE FROM routine_exercises WHERE routine_id = ?", idBytes); err != nil {
		return err
	}
	// 子を再度INSERT
	for _, routineExercise := range routine.Exercises() {
		routineExerciseIDBytes, err := routineExercise.ID().Bytes()
		if err != nil {
			return err
		}
		exerciseIDBytes, err := routineExercise.ExerciseID().Bytes()
		if err != nil {
			return err
		}
		if _, err := q.Exec(insertRoutineExerciseSQL, routineExerciseIDBytes, idBytes, exerciseIDBytes, routineExercise.DisplayOrder().Value(), routineExercise.CreatedAt(), routineExercise.UpdatedAt()); err != nil {
			return err
		}
	}
	return nil
}

func (r *routineRepository) DeleteByID(ctx context.Context, id valueobject.RoutineID) error {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return err
	}
	if _, err := q.Exec("DELETE FROM routines WHERE id = ?", idBytes); err != nil {
		return err
	}
	return nil
}

func toRoutine(row RoutineModel, exercises []*trainingdomain.RoutineExercise) (*trainingdomain.Routine, error) {
	routineID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineID](row.ID)
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
	return trainingdomain.NewRoutine(*routineID, *userID, *name, row.CreatedAt, row.UpdatedAt, exercises), nil
}

func (r *routineRepository) loadExercises(ctx context.Context, routineIDs [][]byte) (map[string][]*trainingdomain.RoutineExercise, error) {
	if len(routineIDs) == 0 {
		return map[string][]*trainingdomain.RoutineExercise{}, nil
	}
	q := dbtx.Querier(ctx, r.dbmap)
	placeholders, args := sqlquery.InPlaceholders(routineIDs)
	var rows []RoutineExerciseModel
	_, err := q.Select(&rows,
		"SELECT id, routine_id, exercise_id, display_order, created_at, updated_at FROM routine_exercises WHERE routine_id IN ("+placeholders+") ORDER BY display_order ASC",
		args...,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return map[string][]*trainingdomain.RoutineExercise{}, nil
	}

	result := make(map[string][]*trainingdomain.RoutineExercise, len(rows))
	for _, row := range rows {
		routineID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineID](row.RoutineID)
		if err != nil {
			return nil, err
		}
		exerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.ExerciseID](row.ExerciseID)
		if err != nil {
			return nil, err
		}
		re, err := toRoutineExercise(row, *routineID, *exerciseID)
		if err != nil {
			return nil, err
		}
		result[routineID.String()] = append(result[routineID.String()], re)
	}
	return result, nil
}

func toRoutineExercise(row RoutineExerciseModel, routineID valueobject.RoutineID, exerciseID valueobject.ExerciseID) (*trainingdomain.RoutineExercise, error) {
	routineExerciseID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.RoutineExerciseID](row.ID)
	if err != nil {
		return nil, err
	}
	displayOrder, err := valueobject.NewNonNegativeInt(int(row.DisplayOrder))
	if err != nil {
		return nil, err
	}
	return trainingdomain.NewRoutineExercise(*routineExerciseID, routineID, exerciseID, *displayOrder, row.CreatedAt, row.UpdatedAt), nil
}
