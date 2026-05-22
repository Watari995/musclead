package valueobject

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

type HashedPassword struct {
	LiteralBase[string]
}

func validateBcryptHash(s any) error {
	if _, err := bcrypt.Cost([]byte(s.(string))); err != nil {
		return errors.New("must be a valid bcrypt hash")
	}
	return nil
}

var hashedPasswordValidationRules = []validation.Rule{
	validation.Required,
	validation.By(validateBcryptHash),
}

func (h HashedPassword) Validate() error {
	return validation.Validate(h.Value(), hashedPasswordValidationRules...)
}

func NewHashedPassword(s string) (*HashedPassword, error) {
	v := HashedPassword{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
