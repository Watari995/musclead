package valueobject

import "github.com/shopspring/decimal"

type DecimalBase struct {
	v decimal.Decimal
}

func (d DecimalBase) Value() decimal.Decimal {
	return d.v
}

func (d DecimalBase) String() string {
	return d.v.String()
}

func (d DecimalBase) Equals(o DecimalBase) bool {
	return d.v.Equal(o.v)
}
