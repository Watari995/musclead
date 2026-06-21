package healthsyncinfra

import "time"

type HealthPlanetTokenModel struct {
	ID           string     `db:"id"`
	UserID       string     `db:"user_id"`
	AccessToken  string     `db:"access_token"`
	RefreshToken string     `db:"refresh_token"`
	ExpiresAt    time.Time  `db:"expires_at"`
	LastSyncedAt *time.Time `db:"last_synced_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
