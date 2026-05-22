package valueobject

import validation "github.com/go-ozzo/ozzo-validation/v4"

type String1000 struct {
	LiteralBase[string]
}

func (s String1000) Validate() error {
	return validation.Validate(s.Value(), validation.Required, validation.RuneLength(1, 1000))
}

func NewString1000(s string) (*String1000, error) {
	v := String1000{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
