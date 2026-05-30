package sessiondomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type TokenSigner interface {
	SignAccessToken(userID valueobject.UserID, expiresAt time.Time) (string, error)
	VerifyAccessToken(token string) (valueobject.UserID, error)
}
