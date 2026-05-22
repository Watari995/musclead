package valueobject

import validation "github.com/go-ozzo/ozzo-validation/v4"

type String50 struct {
	LiteralBase[string]
}

var string50ValidationRules = []validation.Rule{
	validation.Required,
	validation.RuneLength(1, 50),
}

func (s String50) Validate() error {
	return validation.Validate(s.Value(), string50ValidationRules...)
}

func NewString50(s string) (*String50, error) {
	v := String50{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
