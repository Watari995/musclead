package valueobject

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type NonNegativeDecimal struct {
	DecimalBase
}

func validateNonNegativeDecimal(s any) error {
	d, ok := s.(decimal.Decimal)
	if !ok {
		return errors.New("must be a valid decimal")
	}
	if d.LessThan(decimal.Zero) {
		return errors.New("must be greater than or equal to 0")
	}
	return nil
}

func (n NonNegativeDecimal) Validate() error {
	return validation.Validate(n.Value(), validation.By(validateNonNegativeDecimal))
}

func NewNonNegativeDecimal(v decimal.Decimal) (*NonNegativeDecimal, error) {
	n := NonNegativeDecimal{DecimalBase: DecimalBase{v: v}}
	if err := n.Validate(); err != nil {
		return nil, err
	}
	return &n, nil
}
