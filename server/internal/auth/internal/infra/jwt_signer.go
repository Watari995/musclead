package sessioninfra

import (
	"errors"
	"time"

	sessiondomain "github.com/Watari995/musclead/internal/auth/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/golang-jwt/jwt/v5"
)

type accessClaims struct {
	jwt.RegisteredClaims
}

type jwtSigner struct {
	secret []byte
}

func NewJWTSigner(secret string) sessiondomain.TokenSigner {
	return &jwtSigner{secret: []byte(secret)}
}

func (s *jwtSigner) SignAccessToken(userID valueobject.UserID, expiresAt time.Time) (string, error) {
	claims := accessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtSigner) VerifyAccessToken(tokenStr string) (valueobject.UserID, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accessClaims{}, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return valueobject.UserID{}, errors.New("invalid token")
	}
	claims := token.Claims.(*accessClaims)
	userID, err := valueobject.NewPrimaryIDFromString[valueobject.UserID](claims.Subject)
	if err != nil {
		return valueobject.UserID{}, err
	}
	return *userID, nil
}
