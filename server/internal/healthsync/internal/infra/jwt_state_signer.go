package healthsyncinfra

import (
	"errors"
	"time"

	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/golang-jwt/jwt/v5"
)

type stateClaims struct {
	jwt.RegisteredClaims
}

type jwtStateSigner struct {
	secret []byte
}

func NewJWTStateSigner(secret string) healthsyncdomain.StateSigner {
	return &jwtStateSigner{secret: []byte(secret)}
}

func (s *jwtStateSigner) Sign(userID valueobject.UserID) (string, error) {
	claims := stateClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtStateSigner) Verify(state string) (valueobject.UserID, error) {
	token, err := jwt.ParseWithClaims(state, &stateClaims{}, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return valueobject.UserID{}, errors.New("invalid state token")
	}
	claims := token.Claims.(*stateClaims)
	userID, err := valueobject.NewPrimaryIDFromString[valueobject.UserID](claims.Subject)
	if err != nil {
		return valueobject.UserID{}, err
	}
	return *userID, nil
}
