package mealdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type MealTemplate struct {
	id            valueobject.MealTemplateID
	userID        valueobject.UserID
	name          valueobject.String100
	displayOrder  valueobject.NonNegativeInt
	mealType      valueobject.String20
	calories      valueobject.NonNegativeInt
	proteinG      *valueobject.NonNegativeDecimal
	fatG          *valueobject.NonNegativeDecimal
	carbohydrateG *valueobject.NonNegativeDecimal
	createdAt     time.Time
	updatedAt     time.Time
}

func (m *MealTemplate) ID() valueobject.MealTemplateID {
	return m.id
}

func (m *MealTemplate) UserID() valueobject.UserID {
	return m.userID
}

func (m *MealTemplate) Name() valueobject.String100 {
	return m.name
}

func (m *MealTemplate) SetName(name valueobject.String100) {
	m.name = name
	m.updatedAt = time.Now()
}

func (m *MealTemplate) DisplayOrder() valueobject.NonNegativeInt {
	return m.displayOrder
}

func (m *MealTemplate) SetDisplayOrder(displayOrder valueobject.NonNegativeInt) {
	m.displayOrder = displayOrder
	m.updatedAt = time.Now()
}

func (m *MealTemplate) MealType() valueobject.String20 {
	return m.mealType
}

func (m *MealTemplate) SetMealType(mealType valueobject.String20) {
	m.mealType = mealType
	m.updatedAt = time.Now()
}

func (m *MealTemplate) Calories() valueobject.NonNegativeInt {
	return m.calories
}

func (m *MealTemplate) SetCalories(calories valueobject.NonNegativeInt) {
	m.calories = calories
	m.updatedAt = time.Now()
}

func (m *MealTemplate) ProteinG() *valueobject.NonNegativeDecimal {
	return m.proteinG
}

func (m *MealTemplate) SetProteinG(proteinG *valueobject.NonNegativeDecimal) {
	m.proteinG = proteinG
	m.updatedAt = time.Now()
}

func (m *MealTemplate) FatG() *valueobject.NonNegativeDecimal {
	return m.fatG
}

func (m *MealTemplate) SetFatG(fatG *valueobject.NonNegativeDecimal) {
	m.fatG = fatG
	m.updatedAt = time.Now()
}

func (m *MealTemplate) CarbohydrateG() *valueobject.NonNegativeDecimal {
	return m.carbohydrateG
}

func (m *MealTemplate) SetCarbohydrateG(carbohydrateG *valueobject.NonNegativeDecimal) {
	m.carbohydrateG = carbohydrateG
	m.updatedAt = time.Now()
}

func (m *MealTemplate) CreatedAt() time.Time {
	return m.createdAt
}

func (m *MealTemplate) UpdatedAt() time.Time {
	return m.updatedAt
}

type UpdateMealTemplateParams struct {
	Name          valueobject.String100
	MealType      valueobject.String20
	Calories      valueobject.NonNegativeInt
	ProteinG      *valueobject.NonNegativeDecimal
	FatG          *valueobject.NonNegativeDecimal
	CarbohydrateG *valueobject.NonNegativeDecimal
}

func (m *MealTemplate) Update(
	params UpdateMealTemplateParams,
) {
	m.name = params.Name
	m.mealType = params.MealType
	m.calories = params.Calories
	m.proteinG = params.ProteinG
	m.fatG = params.FatG
	m.carbohydrateG = params.CarbohydrateG
	m.updatedAt = time.Now()
}

func CreateMealTemplate(
	userID valueobject.UserID,
	name valueobject.String100,
	displayOrder valueobject.NonNegativeInt,
	mealType valueobject.String20,
	calories valueobject.NonNegativeInt,
	proteinG *valueobject.NonNegativeDecimal,
	fatG *valueobject.NonNegativeDecimal,
	carbohydrateG *valueobject.NonNegativeDecimal,
) *MealTemplate {
	now := time.Now()
	return &MealTemplate{
		id:            valueobject.NewPrimaryID[valueobject.MealTemplateID](),
		userID:        userID,
		name:          name,
		displayOrder:  displayOrder,
		mealType:      mealType,
		calories:      calories,
		proteinG:      proteinG,
		fatG:          fatG,
		carbohydrateG: carbohydrateG,
		createdAt:     now,
		updatedAt:     now,
	}
}

func NewMealTemplate(
	id valueobject.MealTemplateID,
	userID valueobject.UserID,
	name valueobject.String100,
	displayOrder valueobject.NonNegativeInt,
	mealType valueobject.String20,
	calories valueobject.NonNegativeInt,
	proteinG *valueobject.NonNegativeDecimal,
	fatG *valueobject.NonNegativeDecimal,
	carbohydrateG *valueobject.NonNegativeDecimal,
	createdAt time.Time,
	updatedAt time.Time,
) *MealTemplate {
	return &MealTemplate{
		id:            id,
		userID:        userID,
		name:          name,
		displayOrder:  displayOrder,
		mealType:      mealType,
		calories:      calories,
		proteinG:      proteinG,
		fatG:          fatG,
		carbohydrateG: carbohydrateG,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
