package valueobject

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var hexColorRegex = `^#(?:[0-9a-fA-F]{3}){1,2}$`

type ColorHex struct {
	LiteralBase[string]
}

var colorHexValidationRules = []validation.Rule{
	validation.Required,
	validation.Match(regexp.MustCompile(hexColorRegex)),
}

func (c ColorHex) Validate() error {
	return validation.Validate(c.Value(), colorHexValidationRules...)
}

func NewColorHex(s string) (*ColorHex, error) {
	v := ColorHex{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
