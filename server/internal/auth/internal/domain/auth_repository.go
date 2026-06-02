package sessiondomain

import "context"

type SessionRepository interface {
	Save(ctx context.Context, session *Session) error
	FindByRefreshHash(ctx context.Context, refreshHash string) (*Session, error)
}
