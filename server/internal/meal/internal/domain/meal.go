package mealdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type PhotoSpec struct {
	ImagePath    string
	DisplayOrder int
}

type Meal struct {
	id            valueobject.MealID
	userID        valueobject.UserID
	eatenAt       time.Time
	mealType      valueobject.String20
	calories      valueobject.NonNegativeInt
	proteinG      *valueobject.NonNegativeDecimal
	fatG          *valueobject.NonNegativeDecimal
	carbohydrateG *valueobject.NonNegativeDecimal
	memo          *valueobject.String1000
	createdAt     time.Time
	updatedAt     time.Time

	photos []PhotoSpec
}

func (m *Meal) ID() valueobject.MealID {
	return m.id
}

func (m *Meal) UserID() valueobject.UserID {
	return m.userID
}

func (m *Meal) EatenAt() time.Time {
	return m.eatenAt
}

func (m *Meal) MealType() valueobject.String20 {
	return m.mealType
}

func (m *Meal) Calories() valueobject.NonNegativeInt {
	return m.calories
}

func (m *Meal) ProteinG() *valueobject.NonNegativeDecimal {
	return m.proteinG
}

func (m *Meal) FatG() *valueobject.NonNegativeDecimal {
	return m.fatG
}

func (m *Meal) CarbohydrateG() *valueobject.NonNegativeDecimal {
	return m.carbohydrateG
}

func (m *Meal) Memo() *valueobject.String1000 {
	return m.memo
}

func (m *Meal) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Meal) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *Meal) Photos() []PhotoSpec {
	return m.photos
}

func (m *Meal) ReplacePhotos(photos []PhotoSpec) {
	m.photos = photos
	m.updatedAt = time.Now()
}

type UpdateMealParams struct {
	EatenAt       time.Time
	MealType      valueobject.String20
	Calories      valueobject.NonNegativeInt
	ProteinG      *valueobject.NonNegativeDecimal
	FatG          *valueobject.NonNegativeDecimal
	CarbohydrateG *valueobject.NonNegativeDecimal
	Memo          *valueobject.String1000
	Photos        []PhotoSpec
}

func (m *Meal) Update(
	params UpdateMealParams,
) {
	m.eatenAt = params.EatenAt
	m.mealType = params.MealType
	m.calories = params.Calories
	m.proteinG = params.ProteinG
	m.fatG = params.FatG
	m.carbohydrateG = params.CarbohydrateG
	m.memo = params.Memo
	m.photos = params.Photos
	m.updatedAt = time.Now()
}

func CreateMeal(
	userID valueobject.UserID,
	eatenAt time.Time,
	mealType valueobject.String20,
	calories valueobject.NonNegativeInt,
	proteinG *valueobject.NonNegativeDecimal,
	fatG *valueobject.NonNegativeDecimal,
	carbohydrateG *valueobject.NonNegativeDecimal,
	memo *valueobject.String1000,
	photos []PhotoSpec,
) *Meal {
	now := time.Now()
	return &Meal{
		id:            valueobject.NewPrimaryID[valueobject.MealID](),
		userID:        userID,
		eatenAt:       eatenAt,
		mealType:      mealType,
		calories:      calories,
		proteinG:      proteinG,
		fatG:          fatG,
		carbohydrateG: carbohydrateG,
		memo:          memo,
		createdAt:     now,
		updatedAt:     now,
		photos:        photos,
	}
}

func NewMeal(
	id valueobject.MealID,
	userID valueobject.UserID,
	eatenAt time.Time,
	mealType valueobject.String20,
	calories valueobject.NonNegativeInt,
	proteinG *valueobject.NonNegativeDecimal,
	fatG *valueobject.NonNegativeDecimal,
	carbohydrateG *valueobject.NonNegativeDecimal,
	memo *valueobject.String1000,
	createdAt time.Time,
	updatedAt time.Time,
	photos []PhotoSpec,
) *Meal {
	return &Meal{
		id:            id,
		userID:        userID,
		eatenAt:       eatenAt,
		mealType:      mealType,
		calories:      calories,
		proteinG:      proteinG,
		fatG:          fatG,
		carbohydrateG: carbohydrateG,
		memo:          memo,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		photos:        photos,
	}
}
