package valueobject

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Barcode struct {
	LiteralBase[string]
}

// barcodeは13桁以下の数値(8-13桁)
var barcodeValidationRules = []validation.Rule{
	validation.Required,
	validation.RuneLength(8, 13),
	is.Digit,
}

func (b Barcode) Validate() error {
	return validation.Validate(b.Value(), barcodeValidationRules...)
}

func NewBarcode(s string) (*Barcode, error) {
	v := Barcode{LiteralBase: LiteralBase[string]{v: s}}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
