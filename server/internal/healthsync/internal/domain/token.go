package healthsyncdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type Token struct {
	id           valueobject.TokenID
	userID       valueobject.UserID
	accessToken  string
	refreshToken string
	expiresAt    time.Time
	lastSyncedAt *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

func (t *Token) ID() valueobject.TokenID {
	return t.id
}

func (t *Token) UserID() valueobject.UserID {
	return t.userID
}

func (t *Token) AccessToken() string {
	return t.accessToken
}

func (t *Token) RefreshToken() string {
	return t.refreshToken
}

func (t *Token) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t *Token) LastSyncedAt() *time.Time {
	return t.lastSyncedAt
}

func (t *Token) SetLastSyncedAt(lastSyncedAt time.Time) {
	t.lastSyncedAt = &lastSyncedAt
	t.updatedAt = time.Now()
}

func (t *Token) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Token) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Token) UpdateTokens(
	accessToken, refreshToken string,
	expiresAt time.Time,
) {
	t.accessToken = accessToken
	t.refreshToken = refreshToken
	t.expiresAt = expiresAt
	t.updatedAt = time.Now()
}

func CreateToken(
	userID valueobject.UserID,
	accessToken string,
	refreshToken string,
	expiresAt time.Time,
) *Token {
	now := time.Now()
	return &Token{
		id:           valueobject.NewPrimaryID[valueobject.TokenID](),
		userID:       userID,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
		lastSyncedAt: nil,
		createdAt:    now,
		updatedAt:    now,
	}
}

func NewToken(
	id valueobject.TokenID,
	userID valueobject.UserID,
	accessToken string,
	refreshToken string,
	expiresAt time.Time,
	lastSyncedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Token {
	return &Token{
		id:           id,
		userID:       userID,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
		lastSyncedAt: lastSyncedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
