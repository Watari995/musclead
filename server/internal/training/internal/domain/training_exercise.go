package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type TrainingExercise struct {
	id           valueobject.TrainingExerciseID
	trainingID   valueobject.TrainingID
	exerciseID   valueobject.ExerciseID
	displayOrder valueobject.NonNegativeInt
	restSeconds  *valueobject.NonNegativeInt
	memo         *valueobject.String1000
	createdAt    time.Time
	updatedAt    time.Time

	sets []*TrainingSet
}

func (e *TrainingExercise) ID() valueobject.TrainingExerciseID {
	return e.id
}

func (e *TrainingExercise) TrainingID() valueobject.TrainingID {
	return e.trainingID
}

func (e *TrainingExercise) ExerciseID() valueobject.ExerciseID {
	return e.exerciseID
}

func (e *TrainingExercise) DisplayOrder() valueobject.NonNegativeInt {
	return e.displayOrder
}

func (e *TrainingExercise) RestSeconds() *valueobject.NonNegativeInt {
	return e.restSeconds
}

func (e *TrainingExercise) Memo() *valueobject.String1000 {
	return e.memo
}

func (e *TrainingExercise) CreatedAt() time.Time {
	return e.createdAt
}

func (e *TrainingExercise) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *TrainingExercise) Sets() []*TrainingSet {
	return e.sets
}

func CreateTrainingExercise(
	trainingID valueobject.TrainingID,
	exerciseID valueobject.ExerciseID,
	displayOrder valueobject.NonNegativeInt,
	restSeconds *valueobject.NonNegativeInt,
	memo *valueobject.String1000,
	sets []*TrainingSet,
) *TrainingExercise {
	now := time.Now()
	return &TrainingExercise{
		id:           valueobject.NewPrimaryID[valueobject.TrainingExerciseID](),
		trainingID:   trainingID,
		exerciseID:   exerciseID,
		displayOrder: displayOrder,
		restSeconds:  restSeconds,
		memo:         memo,
		createdAt:    now,
		updatedAt:    now,
		sets:         sets,
	}
}

func NewTrainingExercise(
	id valueobject.TrainingExerciseID,
	trainingID valueobject.TrainingID,
	exerciseID valueobject.ExerciseID,
	displayOrder valueobject.NonNegativeInt,
	restSeconds *valueobject.NonNegativeInt,
	memo *valueobject.String1000,
	createdAt time.Time,
	updatedAt time.Time,
	sets []*TrainingSet,
) *TrainingExercise {
	return &TrainingExercise{
		id:           id,
		trainingID:   trainingID,
		exerciseID:   exerciseID,
		displayOrder: displayOrder,
		restSeconds:  restSeconds,
		memo:         memo,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		sets:         sets,
	}
}
