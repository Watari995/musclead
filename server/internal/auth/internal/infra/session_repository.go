package sessioninfra

import (
	"context"
	"database/sql"
	"errors"

	sessiondomain "github.com/Watari995/musclead/internal/auth/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type sessionRepository struct {
	dbmap *gorp.DbMap
}

func (r *sessionRepository) Save(ctx context.Context, session *sessiondomain.Session) error {
	q := dbtx.Querier(ctx, r.dbmap)
	bytes, err := session.ID().Bytes()
	if err != nil {
		return err
	}
	userIDBytes, err := session.UserID().Bytes()
	if err != nil {
		return err
	}
	params := []interface{}{
		bytes,
		userIDBytes,
		session.RefreshHash(),
		session.UserAgent(),
		session.IPAddress(),
		session.ExpiresAt(),
		sqlconv.ToNullTime(session.RevokedAt()),
		session.CreatedAt(),
	}
	_, err = q.Exec(`
INSERT INTO sessions (id, user_id, refresh_hash, user_agent, ip_address, expires_at, revoked_at, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    refresh_hash = VALUES(refresh_hash),
    expires_at = VALUES(expires_at),
    revoked_at = VALUES(revoked_at)`,
		params...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) FindByRefreshHash(ctx context.Context, refreshHash string) (*sessiondomain.Session, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var row SessionModel
	err := q.SelectOne(&row, "SELECT id, user_id, refresh_hash, user_agent, ip_address, expires_at, revoked_at, created_at FROM sessions WHERE refresh_hash = ?", refreshHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toSession(row)
}

func toSession(row SessionModel) (*sessiondomain.Session, error) {
	sessionID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.SessionID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	return sessiondomain.NewSession(*sessionID, *userID, row.RefreshHash, row.UserAgent, row.IPAddress, row.ExpiresAt, sqlconv.FromNullTime(row.RevokedAt), row.CreatedAt), nil
}

func NewSessionRepository(dbmap *gorp.DbMap) sessiondomain.SessionRepository {
	return &sessionRepository{dbmap: dbmap}
}
