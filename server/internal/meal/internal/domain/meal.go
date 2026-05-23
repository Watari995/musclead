package mealdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type PhotoData struct {
	ImagePath    string
	DisplayOrder int
}

type Meal struct {
	id            valueobject.MealID
	userId        valueobject.UserID
	eatenAt       time.Time
	mealType      valueobject.String20
	calories      valueobject.NonNegativeInt
	proteinG      valueobject.NonNegativeInt
	fatG          valueobject.NonNegativeInt
	carbohydrateG valueobject.NonNegativeInt
	memo          *valueobject.String1000
	createdAt     time.Time
	updatedAt     time.Time

	photos []PhotoData
}

func (m *Meal) ID() valueobject.MealID {
	return m.id
}

func (m *Meal) UserID() valueobject.UserID {
	return m.userId
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

func (m *Meal) ProteinG() valueobject.NonNegativeInt {
	return m.proteinG
}

func (m *Meal) FatG() valueobject.NonNegativeInt {
	return m.fatG
}

func (m *Meal) CarbohydrateG() valueobject.NonNegativeInt {
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

func (m *Meal) Photos() []PhotoData {
	return m.photos
}

func (m *Meal) ReplacePhotos(photos []PhotoData) {
	m.photos = photos
	m.updatedAt = time.Now()
}

func CreateMeal(userId valueobject.UserID, eatenAt time.Time, mealType valueobject.String20, calories valueobject.NonNegativeInt, proteinG valueobject.NonNegativeInt, fatG valueobject.NonNegativeInt, carbohydrateG valueobject.NonNegativeInt, memo *valueobject.String1000, photos []PhotoData) *Meal {
	now := time.Now()
	return &Meal{
		id:            valueobject.NewPrimaryId[valueobject.MealID](),
		userId:        userId,
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

func NewMeal(id valueobject.MealID, userId valueobject.UserID, eatenAt time.Time, mealType valueobject.String20, calories valueobject.NonNegativeInt, proteinG valueobject.NonNegativeInt, fatG valueobject.NonNegativeInt, carbohydrateG valueobject.NonNegativeInt, memo *valueobject.String1000, createdAt time.Time, updatedAt time.Time, photos []PhotoData) *Meal {
	return &Meal{
		id:            id,
		userId:        userId,
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
