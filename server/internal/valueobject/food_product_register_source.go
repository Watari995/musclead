package valueobject

import "errors"

type FoodProductRegisterSourceCode string

const (
	FoodProductRegisterSourceOpenFoodFacts FoodProductRegisterSourceCode = "open_food_facts"
	FoodProductRegisterSourceUser          FoodProductRegisterSourceCode = "user"
)

var ErrInvalidFoodProductRegisterSourceCode = errors.New("invalid food product register source code")

type FoodProductRegisterSource struct {
	LiteralBase[string]
}

func NewFoodProductRegisterSourceFromString(s string) (*FoodProductRegisterSource, error) {
	switch FoodProductRegisterSourceCode(s) {
	case FoodProductRegisterSourceOpenFoodFacts, FoodProductRegisterSourceUser:
		return &FoodProductRegisterSource{LiteralBase: LiteralBase[string]{v: s}}, nil
	default:
		return nil, ErrInvalidFoodProductRegisterSourceCode
	}
}

func NewFoodProductRegisterSourceFromCode(c FoodProductRegisterSourceCode) (*FoodProductRegisterSource, error) {
	return &FoodProductRegisterSource{LiteralBase: LiteralBase[string]{v: string(c)}}, nil
}
