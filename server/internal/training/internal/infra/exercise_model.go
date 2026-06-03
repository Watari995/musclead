package traininginfra

import "time"

type ExerciseModel struct {
	ID        []byte    `db:"id"`
	UserID    []byte    `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
