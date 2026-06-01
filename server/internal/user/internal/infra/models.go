package userinfra

import (
	"database/sql"
	"time"
)

type UserModel struct {
	ID           []byte       `db:"id"`
	Name         string       `db:"name"`
	Email        string       `db:"email"`
	PasswordHash string       `db:"password_hash"`
	Birthday     sql.NullTime `db:"birthday"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}
