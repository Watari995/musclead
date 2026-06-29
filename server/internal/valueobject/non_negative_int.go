package valueobject

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type NonNegativeInt struct {
	LiteralBase[int]
}

var nonNegativeIntValidationRules = []validation.Rule{
	validation.Min(0),
}

func (n NonNegativeInt) Validate() error {
	return validation.Validate(n.Value(), nonNegativeIntValidationRules...)
}

func NewNonNegativeInt(v int) (*NonNegativeInt, error) {
	n := NonNegativeInt{LiteralBase: LiteralBase[int]{v: v}}
	if err := n.Validate(); err != nil {
		return nil, err
	}
	return &n, nil
}

func (n NonNegativeInt) Add(other NonNegativeInt) NonNegativeInt {
	return NonNegativeInt{LiteralBase: LiteralBase[int]{v: n.Value() + other.Value()}}
}
