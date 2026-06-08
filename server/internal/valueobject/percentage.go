package valueobject

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type Percentage struct {
	DecimalBase
}

func validatePercentage(s any) error {
	d, ok := s.(decimal.Decimal)
	if !ok {
		return errors.New("must be a valid decimal")
	}
	if d.LessThan(decimal.Zero) {
		return errors.New("must be greater than or equal to 0")
	}
	if d.GreaterThan(decimal.NewFromInt(100)) {
		return errors.New("must be less than or equal to 100")
	}
	return nil
}

func (p Percentage) Validate() error {
	return validation.Validate(p.Value(), validation.By(validatePercentage))
}

func NewPercentage(v decimal.Decimal) (*Percentage, error) {
	p := Percentage{DecimalBase: DecimalBase{v: v}}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func NewPercentageFromString(s string) (*Percentage, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return NewPercentage(d)
}
