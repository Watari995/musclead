package valueobject

import validation "github.com/go-ozzo/ozzo-validation/v4"

type String100 struct {
	LiteralBase[string]
}

func (s String100) Validate() error {
	return validation.Validate(s.Value(), validation.Required, validation.RuneLength(1, 100))
}

func NewString100(s string) (*String100, error) {
	v := String100{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
