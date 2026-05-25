package userinfra

import (
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordHasher struct {
	cost int
}

func (h *bcryptPasswordHasher) Hash(rawPassword string) (*valueobject.HashedPassword, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(rawPassword), h.cost)
	if err != nil {
		return nil, err
	}
	return valueobject.NewHashedPassword(string(b))
}

func (h *bcryptPasswordHasher) Compare(rawPassword string, hash *valueobject.HashedPassword) error {
	return bcrypt.CompareHashAndPassword([]byte(hash.Value()), []byte(rawPassword))
}

func NewBcryptPasswordHasher() userdomain.PasswordHasher {
	return &bcryptPasswordHasher{cost: bcrypt.DefaultCost}
}
