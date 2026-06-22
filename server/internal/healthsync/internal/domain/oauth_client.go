package healthsyncdomain

import (
	"context"
	"time"
)

// TokenExchanger は OAuth 認可コードをトークンに交換する。
// connect_health_planet use case が依存する。
type TokenExchanger interface {
	ExchangeCode(ctx context.Context, code, redirectURI string) (accessToken, refreshToken string, expiresAt time.Time, err error)
}

// TokenRefresher はアクセストークンをリフレッシュする。
// sync_weights use case が依存する。
type TokenRefresher interface {
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, expiresAt time.Time, err error)
}
