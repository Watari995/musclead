package userinfra

import "time"

type UserPreferencesModel struct {
	ID        []byte    `db:"id"`
	UserID    []byte    `db:"user_id"`
	Theme     string    `db:"theme"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
