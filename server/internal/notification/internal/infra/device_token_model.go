package notificationinfra

import "time"

type DeviceTokenModel struct {
	ID        []byte    `db:"id"`
	UserID    []byte    `db:"user_id"`
	Token     string    `db:"token"`
	Platform  string    `db:"platform"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
