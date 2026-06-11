package valueobject

import "errors"

type CurrencyCode string

const (
	CurrencyJPY CurrencyCode = "JPY"
	CurrencyUSD CurrencyCode = "USD"
	CurrencyEUR CurrencyCode = "EUR"
	CurrencyGBP CurrencyCode = "GBP"
	CurrencyCAD CurrencyCode = "CAD"
	CurrencyAUD CurrencyCode = "AUD"
	CurrencyNZD CurrencyCode = "NZD"
	CurrencyCHF CurrencyCode = "CHF"
	CurrencyCNY CurrencyCode = "CNY"
)

var ErrInvalidCurrency = errors.New("invalid currency")

type Currency struct {
	LiteralBase[string]
}

func NewCurrencyFromString(s string) (*Currency, error) {
	switch s {
	case string(CurrencyJPY), string(CurrencyUSD), string(CurrencyEUR), string(CurrencyGBP), string(CurrencyCAD), string(CurrencyAUD), string(CurrencyNZD), string(CurrencyCHF), string(CurrencyCNY):
		return &Currency{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidCurrency
	}
}

func NewCurrencyFromCode(c CurrencyCode) Currency {
	return Currency{LiteralBase: LiteralBase[string]{v: string(c)}}
}
