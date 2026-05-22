package valueobject

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Email struct {
	LiteralBase[string]
}

var emailValidationRules = []validation.Rule{
	validation.Required,
	validation.RuneLength(1, 255),
	is.Email,
}

func (e Email) Validate() error {
	return validation.Validate(e.Value(), emailValidationRules...)
}

func NewEmail(s string) (*Email, error) {
	v := Email{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
