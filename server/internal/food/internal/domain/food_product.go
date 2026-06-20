package fooddomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// FoodProduct は食品マスタのエンティティ。
type FoodProduct struct {
	id             valueobject.FoodProductID
	barcode        *valueobject.Barcode
	name           valueobject.String100
	calories       valueobject.NonNegativeInt
	proteinG       *valueobject.NonNegativeDecimal
	fatG           *valueobject.NonNegativeDecimal
	carbohydrateG  *valueobject.NonNegativeDecimal
	registerSource valueobject.FoodProductRegisterSource
	createdAt      time.Time
	updatedAt      time.Time
}

func (f *FoodProduct) ID() valueobject.FoodProductID {
	return f.id
}

func (f *FoodProduct) Barcode() *valueobject.Barcode {
	return f.barcode
}

func (f *FoodProduct) Name() valueobject.String100 {
	return f.name
}

func (f *FoodProduct) Calories() valueobject.NonNegativeInt {
	return f.calories
}

func (f *FoodProduct) ProteinG() *valueobject.NonNegativeDecimal {
	return f.proteinG
}

func (f *FoodProduct) FatG() *valueobject.NonNegativeDecimal {
	return f.fatG
}

func (f *FoodProduct) CarbohydrateG() *valueobject.NonNegativeDecimal {
	return f.carbohydrateG
}

func (f *FoodProduct) RegisterSource() valueobject.FoodProductRegisterSource {
	return f.registerSource
}

func (f *FoodProduct) CreatedAt() time.Time {
	return f.createdAt
}

func (f *FoodProduct) UpdatedAt() time.Time {
	return f.updatedAt
}

func CreateFoodProduct(
	barcode *valueobject.Barcode,
	name valueobject.String100,
	calories valueobject.NonNegativeInt,
	proteinG *valueobject.NonNegativeDecimal,
	fatG *valueobject.NonNegativeDecimal,
	carbohydrateG *valueobject.NonNegativeDecimal,
	registerSource valueobject.FoodProductRegisterSource,
) *FoodProduct {
	now := time.Now()
	return &FoodProduct{
		id:             valueobject.NewPrimaryID[valueobject.FoodProductID](),
		barcode:        barcode,
		name:           name,
		calories:       calories,
		proteinG:       proteinG,
		fatG:           fatG,
		carbohydrateG:  carbohydrateG,
		registerSource: registerSource,
		createdAt:      now,
		updatedAt:      now,
	}
}

func NewFoodProduct(
	id valueobject.FoodProductID,
	barcode *valueobject.Barcode,
	name valueobject.String100,
	calories valueobject.NonNegativeInt,
	proteinG *valueobject.NonNegativeDecimal,
	fatG *valueobject.NonNegativeDecimal,
	carbohydrateG *valueobject.NonNegativeDecimal,
	registerSource valueobject.FoodProductRegisterSource,
	createdAt time.Time,
	updatedAt time.Time,
) *FoodProduct {
	return &FoodProduct{
		id:             id,
		barcode:        barcode,
		name:           name,
		calories:       calories,
		proteinG:       proteinG,
		fatG:           fatG,
		carbohydrateG:  carbohydrateG,
		registerSource: registerSource,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}
