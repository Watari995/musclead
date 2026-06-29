package valueobject

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type WeightChangeKg struct {
	DecimalBase
}

func validateWeightChangeKg(v any) error {
	s, ok := v.(decimal.Decimal)
	if !ok {
		return errors.New("must be a valid decimal")
	}
	if s.LessThan(decimal.NewFromInt(-100)) {
		return errors.New("must be greater than or equal to -100")
	}
	if s.GreaterThan(decimal.NewFromInt(100)) {
		return errors.New("must be less than or equal to 100")
	}
	return nil
}

func (w WeightChangeKg) Validate() error {
	return validation.Validate(w.Value(), validation.By(validateWeightChangeKg))
}

func NewWeightChangeKgFromDecimal(v decimal.Decimal) (*WeightChangeKg, error) {
	w := WeightChangeKg{DecimalBase: DecimalBase{v: v}}
	if err := w.Validate(); err != nil {
		return nil, err
	}
	return &w, nil
}
