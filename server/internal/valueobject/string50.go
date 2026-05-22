package valueobject

import validation "github.com/go-ozzo/ozzo-validation/v4"

type String50 struct {
	LiteralBase[string]
}

func (s String50) Validate() error {
	return validation.Validate(s.Value(), validation.Required, validation.RuneLength(1, 50))
}

func NewString50(s string) (*String50, error) {
	v := String50{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
