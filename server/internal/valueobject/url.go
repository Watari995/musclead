package valueobject

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type URL struct {
	LiteralBase[string]
}

var urlValidationRules = []validation.Rule{
	validation.Required,
	validation.RuneLength(1, 2048),
	is.URL,
}

func (u URL) Validate() error {
	return validation.Validate(u.Value(), urlValidationRules...)
}

func NewURL(s string) (*URL, error) {
	v := URL{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
