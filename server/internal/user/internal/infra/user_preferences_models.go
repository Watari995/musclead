package userinfra

import "time"

type UserPreferencesModel struct {
	ID            []byte    `db:"id"`
	UserID        []byte    `db:"user_id"`
	Theme         string    `db:"theme"`
	MealColor     string    `db:"meal_color"`
	TrainingColor string    `db:"training_color"`
	WeightColor   string    `db:"weight_color"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
