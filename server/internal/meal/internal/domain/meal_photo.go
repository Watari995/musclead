package mealdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type MealPhoto struct {
	id           valueobject.MealPhotoID
	mealID       valueobject.MealID
	imagePath    string
	displayOrder int
	createdAt    time.Time
}

func (mp *MealPhoto) ID() valueobject.MealPhotoID {
	return mp.id
}

func (mp *MealPhoto) MealID() valueobject.MealID {
	return mp.mealID
}

func (mp *MealPhoto) ImagePath() string {
	return mp.imagePath
}

func (mp *MealPhoto) DisplayOrder() int {
	return mp.displayOrder
}

func (mp *MealPhoto) CreatedAt() time.Time {
	return mp.createdAt
}

func CreateMealPhoto(mealID valueobject.MealID, imagePath string, displayOrder int) *MealPhoto {
	return &MealPhoto{
		id:           valueobject.NewPrimaryId[valueobject.MealPhotoID](),
		mealID:       mealID,
		imagePath:    imagePath,
		displayOrder: displayOrder,
		createdAt:    time.Now(),
	}
}

func NewMealPhoto(id valueobject.MealPhotoID, mealID valueobject.MealID, imagePath string, displayOrder int, createdAt time.Time) *MealPhoto {
	return &MealPhoto{
		id:           id,
		mealID:       mealID,
		imagePath:    imagePath,
		displayOrder: displayOrder,
		createdAt:    createdAt,
	}
}
