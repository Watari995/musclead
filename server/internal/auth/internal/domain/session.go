package sessiondomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type Session struct {
	id          valueobject.SessionID
	userID      valueobject.UserID
	refreshHash string
	userAgent   string
	ipAddress   string
	expiresAt   time.Time
	revokedAt   *time.Time
	createdAt   time.Time
}

func (s *Session) ID() valueobject.SessionID {
	return s.id
}

func (s *Session) UserID() valueobject.UserID {
	return s.userID
}

func (s *Session) RefreshHash() string {
	return s.refreshHash
}

func (s *Session) UserAgent() string {
	return s.userAgent
}

func (s *Session) IPAddress() string {
	return s.ipAddress
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) IsActive() bool {
	return s.revokedAt == nil && s.expiresAt.After(time.Now())
}

func (s *Session) RevokedAt() *time.Time {
	return s.revokedAt
}

func (s *Session) Revoke() {
	if s.revokedAt != nil {
		return
	}
	now := time.Now()
	s.revokedAt = &now
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func CreateSession(
	userID valueobject.UserID,
	refreshHash string,
	userAgent string,
	ipAddress string,
	expiresAt time.Time,
) *Session {
	return &Session{
		id:          valueobject.NewPrimaryID[valueobject.SessionID](),
		userID:      userID,
		refreshHash: refreshHash,
		userAgent:   userAgent,
		ipAddress:   ipAddress,
		expiresAt:   expiresAt,
		revokedAt:   nil,
		createdAt:   time.Now(),
	}
}

func NewSession(
	id valueobject.SessionID,
	userID valueobject.UserID,
	refreshHash string,
	userAgent string,
	ipAddress string,
	expiresAt time.Time,
	revokedAt *time.Time,
	createdAt time.Time,
) *Session {
	return &Session{
		id:          id,
		userID:      userID,
		refreshHash: refreshHash,
		userAgent:   userAgent,
		ipAddress:   ipAddress,
		expiresAt:   expiresAt,
		revokedAt:   revokedAt,
		createdAt:   createdAt,
	}
}
