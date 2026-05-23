package userdomain

import "github.com/Watari995/musclead/internal/valueobject"

type PasswordHasher interface {
	Hash(rawPassword string) (*valueobject.HashedPassword, error)
	Compare(rawPassword string, hash *valueobject.HashedPassword) error
}
