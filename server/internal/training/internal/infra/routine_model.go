package traininginfra

import "time"

type RoutineModel struct {
	ID           []byte    `db:"id"`
	UserID       []byte    `db:"user_id"`
	Name         string    `db:"name"`
	DisplayOrder int32     `db:"display_order"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type RoutineExerciseModel struct {
	ID           []byte    `db:"id"`
	RoutineID    []byte    `db:"routine_id"`
	ExerciseID   []byte    `db:"exercise_id"`
	DisplayOrder int32     `db:"display_order"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
