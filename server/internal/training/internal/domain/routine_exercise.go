package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type RoutineExercise struct {
	id           valueobject.RoutineExerciseID
	routineID    valueobject.RoutineID
	exerciseID   valueobject.ExerciseID
	displayOrder valueobject.NonNegativeInt
	createdAt    time.Time
	updatedAt    time.Time
}

func (e *RoutineExercise) ID() valueobject.RoutineExerciseID {
	return e.id
}

func (e *RoutineExercise) RoutineID() valueobject.RoutineID {
	return e.routineID
}

func (e *RoutineExercise) ExerciseID() valueobject.ExerciseID {
	return e.exerciseID
}

func (e *RoutineExercise) DisplayOrder() valueobject.NonNegativeInt {
	return e.displayOrder
}

func (e *RoutineExercise) CreatedAt() time.Time {
	return e.createdAt
}

func (e *RoutineExercise) UpdatedAt() time.Time {
	return e.updatedAt
}

func CreateRoutineExercise(
	routineID valueobject.RoutineID,
	exerciseID valueobject.ExerciseID,
	displayOrder valueobject.NonNegativeInt,
) *RoutineExercise {
	now := time.Now()
	return &RoutineExercise{
		id:           valueobject.NewPrimaryID[valueobject.RoutineExerciseID](),
		routineID:    routineID,
		exerciseID:   exerciseID,
		displayOrder: displayOrder,
		createdAt:    now,
		updatedAt:    now,
	}
}

func NewRoutineExercise(
	id valueobject.RoutineExerciseID,
	routineID valueobject.RoutineID,
	exerciseID valueobject.ExerciseID,
	displayOrder valueobject.NonNegativeInt,
	createdAt time.Time,
	updatedAt time.Time,
) *RoutineExercise {
	return &RoutineExercise{
		id:           id,
		routineID:    routineID,
		exerciseID:   exerciseID,
		displayOrder: displayOrder,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
