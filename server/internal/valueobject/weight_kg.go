package valueobject

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type WeightKg struct {
	DecimalBase
}

// decimalはstructなので、通常のvalidation ruleではなくてvalidation.By(validateWeightKg)を使う。
func validateWeightKg(s any) error {
	// cast to decimal
	d, ok := s.(decimal.Decimal)
	if !ok {
		return errors.New("must be a valid decimal")
	}
	if d.LessThanOrEqual(decimal.Zero) {
		return errors.New("must be greater than 0")
	}
	if d.GreaterThanOrEqual(decimal.NewFromInt(1000)) {
		return errors.New("must be less than 1000")
	}
	return nil
}

func (w WeightKg) Validate() error {
	return validation.Validate(w.Value(), validation.By(validateWeightKg))
}

func NewWeightKg(v decimal.Decimal) (*WeightKg, error) {
	w := WeightKg{DecimalBase: DecimalBase{v: v}}
	if err := w.Validate(); err != nil {
		return nil, err
	}
	return &w, nil
}
