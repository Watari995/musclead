package healthsyncinfra

import (
	"database/sql"
	"time"
)

type HealthPlanetTokenModel struct {
	ID           []byte       `db:"id"`
	UserID       []byte       `db:"user_id"`
	AccessToken  string       `db:"access_token"`
	RefreshToken string       `db:"refresh_token"`
	ExpiresAt    time.Time    `db:"expires_at"`
	LastSyncedAt sql.NullTime `db:"last_synced_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}
