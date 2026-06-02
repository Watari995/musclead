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

// NewNonNegativeDecimalFromString は数値の文字列表現("12.5" 等) から VO を生成する。
// HTTP 入力 / DB の DECIMAL 列(string で返る)等、 string 起点の変換で共通利用する。
func NewNonNegativeDecimalFromString(s string) (*NonNegativeDecimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return NewNonNegativeDecimal(d)
}

// parse optional non negative decimal for handler parse to vo
func ParseOptionalNonNegativeDecimal(f *float64) (*NonNegativeDecimal, error) {
	if f == nil {
		return nil, nil
	}
	d, err := NewNonNegativeDecimal(decimal.NewFromFloat(*f))
	if err != nil {
		return nil, err
	}
	return d, nil
}
