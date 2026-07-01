package userinfra

import (
	"database/sql"
	"time"
)

type UserWeeklyGoalModel struct {
	ID             []byte         `db:"id"`
	UserID         []byte         `db:"user_id"`
	TrainingCount  sql.NullInt32  `db:"training_count"`
	CalorieAverage sql.NullInt32  `db:"calorie_average"`
	WeightChangeKg sql.NullString `db:"weight_change_kg"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}
