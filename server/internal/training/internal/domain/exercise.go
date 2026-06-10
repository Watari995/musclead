package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type Exercise struct {
	id           valueobject.ExerciseID
	userID       valueobject.UserID
	name         valueobject.String50
	displayOrder valueobject.NonNegativeInt
	createdAt    time.Time
	updatedAt    time.Time
}

func (e *Exercise) ID() valueobject.ExerciseID {
	return e.id
}

func (e *Exercise) UserID() valueobject.UserID {
	return e.userID
}

func (e *Exercise) Name() valueobject.String50 {
	return e.name
}

func (e *Exercise) SetName(name valueobject.String50) {
	e.name = name
	e.updatedAt = time.Now()
}

func (e *Exercise) DisplayOrder() valueobject.NonNegativeInt {
	return e.displayOrder
}

func (e *Exercise) SetDisplayOrder(displayOrder valueobject.NonNegativeInt) {
	e.displayOrder = displayOrder
	e.updatedAt = time.Now()
}

func (e *Exercise) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Exercise) UpdatedAt() time.Time {
	return e.updatedAt
}

func CreateExercise(
	userID valueobject.UserID,
	name valueobject.String50,
	displayOrder valueobject.NonNegativeInt,
) *Exercise {
	now := time.Now()
	return &Exercise{
		id:           valueobject.NewPrimaryID[valueobject.ExerciseID](),
		userID:       userID,
		name:         name,
		displayOrder: displayOrder,
		createdAt:    now,
		updatedAt:    now,
	}
}

func NewExercise(
	id valueobject.ExerciseID,
	userID valueobject.UserID,
	name valueobject.String50,
	displayOrder valueobject.NonNegativeInt,
	createdAt time.Time,
	updatedAt time.Time,
) *Exercise {
	return &Exercise{
		id:           id,
		userID:       userID,
		name:         name,
		displayOrder: displayOrder,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
