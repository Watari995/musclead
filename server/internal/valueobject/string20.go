package valueobject

import validation "github.com/go-ozzo/ozzo-validation/v4"

type String20 struct {
	LiteralBase[string]
}

var string20ValidationRules = []validation.Rule{
	validation.Required,
	validation.RuneLength(1, 20),
}

func (s String20) Validate() error {
	return validation.Validate(s.Value(), string20ValidationRules...)
}

func NewString20(s string) (*String20, error) {
	v := String20{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
