package weightinfra

import (
	"database/sql"
	"time"
)

type WeightModel struct {
	ID                []byte         `db:"id"`
	UserID            []byte         `db:"user_id"`
	WeightKg          string         `db:"weight_kg"`
	BodyFatPercentage sql.NullString `db:"body_fat_percentage"`
	SkeletalMuscleKg  sql.NullString `db:"skeletal_muscle_kg"`
	MeasuredAt        time.Time      `db:"measured_at"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}
