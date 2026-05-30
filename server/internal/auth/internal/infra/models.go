package sessioninfra

import (
	"database/sql"
	"time"
)

type sessionModel struct {
	ID          []byte       `db:"id"`
	UserID      []byte       `db:"user_id"`
	RefreshHash string       `db:"refresh_hash"`
	UserAgent   string       `db:"user_agent"`
	IPAddress   string       `db:"ip_address"`
	ExpiresAt   time.Time    `db:"expires_at"`
	RevokedAt   sql.NullTime `db:"revoked_at"`
	CreatedAt   time.Time    `db:"created_at"`
}
