package traininginfra

import (
	"database/sql"
	"time"
)

type TrainingModel struct {
	ID        []byte       `db:"id"`
	UserID    []byte       `db:"user_id"`
	StartedAt time.Time    `db:"started_at"`
	EndedAt   sql.NullTime `db:"ended_at"`
	Memo      sql.NullString `db:"memo"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}

type TrainingExerciseModel struct {
	ID           []byte         `db:"id"`
	TrainingID   []byte         `db:"training_id"`
	Name         string         `db:"name"`
	DisplayOrder int32          `db:"display_order"`
	RestSeconds  sql.NullInt32  `db:"rest_seconds"`
	Memo         sql.NullString `db:"memo"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type TrainingSetModel struct {
	ID                 []byte         `db:"id"`
	TrainingExerciseID []byte         `db:"training_exercise_id"`
	SetNumber          int32          `db:"set_number"`
	WeightKg           string         `db:"weight_kg"`
	Reps               int32          `db:"reps"`
	RestSeconds        sql.NullInt32  `db:"rest_seconds"`
	Memo               sql.NullString `db:"memo"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}
