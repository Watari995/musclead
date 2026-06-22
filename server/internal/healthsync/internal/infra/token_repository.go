package healthsyncinfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-gorp/gorp/v3"

	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
)

type tokenRepository struct {
	dbmap *gorp.DbMap
}

func NewTokenRepository(dbmap *gorp.DbMap) healthsyncdomain.TokenRepository {
	return &tokenRepository{dbmap: dbmap}
}

const upsertTokenSQL = `
INSERT INTO healthplanet_tokens (id, user_id, access_token, refresh_token, expires_at, last_synced_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	access_token   = VALUES(access_token),
	refresh_token  = VALUES(refresh_token),
	expires_at     = VALUES(expires_at),
	last_synced_at = VALUES(last_synced_at),
	updated_at     = VALUES(updated_at)
`

func (r *tokenRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*healthsyncdomain.Token, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}

	var row HealthPlanetTokenModel
	err = q.SelectOne(&row,
		"SELECT id, user_id, access_token, refresh_token, expires_at, last_synced_at, created_at, updated_at FROM healthplanet_tokens WHERE user_id = ?",
		idBytes,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toEntity(row)
}

func (r *tokenRepository) FindAllActive(ctx context.Context) ([]*healthsyncdomain.Token, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	var rows []HealthPlanetTokenModel
	if _, err := q.Select(&rows, "SELECT id, user_id, access_token, refresh_token, expires_at, last_synced_at, created_at, updated_at FROM healthplanet_tokens"); err != nil {
		return nil, err
	}
	result := make([]*healthsyncdomain.Token, len(rows))
	for i, row := range rows {
		token, err := toEntity(row)
		if err != nil {
			return nil, err
		}
		result[i] = token
	}
	return result, nil
}

func (r *tokenRepository) Save(ctx context.Context, token *healthsyncdomain.Token) error {
	q := dbtx.Querier(ctx, r.dbmap)

	params, err := buildUpsertTokenParams(token)
	if err != nil {
		return err
	}

	if _, err := q.Exec(upsertTokenSQL, params...); err != nil {
		return err
	}
	return nil
}

func buildUpsertTokenParams(token *healthsyncdomain.Token) ([]any, error) {
	idBytes, err := token.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := token.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	var lastSyncedAt sql.NullTime
	if token.LastSyncedAt() != nil {
		lastSyncedAt = sqlconv.ToNullTime(token.LastSyncedAt())
	}
	return []any{
		idBytes,
		userIDBytes,
		token.AccessToken(),
		token.RefreshToken(),
		token.ExpiresAt(),
		lastSyncedAt,
		token.CreatedAt(),
		token.UpdatedAt(),
	}, nil
}

func toEntity(row HealthPlanetTokenModel) (*healthsyncdomain.Token, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.TokenID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	lastSyncedAt := sqlconv.FromNullTime(row.LastSyncedAt)
	return healthsyncdomain.NewToken(
		*id,
		*userID,
		row.AccessToken,
		row.RefreshToken,
		row.ExpiresAt,
		lastSyncedAt,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}
