package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type TrainingSet struct {
	id                 valueobject.SetID
	trainingExerciseID valueobject.ExerciseID
	setNumber          valueobject.NonNegativeInt
	weightKg           valueobject.NonNegativeDecimal
	reps               valueobject.NonNegativeInt
	restSeconds        *valueobject.NonNegativeInt
	memo               *valueobject.String1000
	createdAt          time.Time
	updatedAt          time.Time
}

func (s *TrainingSet) ID() valueobject.SetID {
	return s.id
}

func (s *TrainingSet) TrainingExerciseID() valueobject.ExerciseID {
	return s.trainingExerciseID
}

func (s *TrainingSet) SetNumber() valueobject.NonNegativeInt {
	return s.setNumber
}

func (s *TrainingSet) WeightKg() valueobject.NonNegativeDecimal {
	return s.weightKg
}

func (s *TrainingSet) Reps() valueobject.NonNegativeInt {
	return s.reps
}

func (s *TrainingSet) RestSeconds() *valueobject.NonNegativeInt {
	return s.restSeconds
}

func (s *TrainingSet) Memo() *valueobject.String1000 {
	return s.memo
}

func (s *TrainingSet) CreatedAt() time.Time {
	return s.createdAt
}

func (s *TrainingSet) UpdatedAt() time.Time {
	return s.updatedAt
}

func CreateTrainingSet(
	trainingExerciseID valueobject.ExerciseID,
	setNumber valueobject.NonNegativeInt,
	weightKg valueobject.NonNegativeDecimal,
	reps valueobject.NonNegativeInt,
	restSeconds *valueobject.NonNegativeInt,
	memo *valueobject.String1000,
) *TrainingSet {
	now := time.Now()
	return &TrainingSet{
		id:                 valueobject.NewPrimaryID[valueobject.SetID](),
		trainingExerciseID: trainingExerciseID,
		setNumber:          setNumber,
		weightKg:           weightKg,
		reps:               reps,
		restSeconds:        restSeconds,
		memo:               memo,
		createdAt:          now,
		updatedAt:          now,
	}
}

func NewTrainingSet(
	id valueobject.SetID,
	trainingExerciseID valueobject.ExerciseID,
	setNumber valueobject.NonNegativeInt,
	weightKg valueobject.NonNegativeDecimal,
	reps valueobject.NonNegativeInt,
	restSeconds *valueobject.NonNegativeInt,
	memo *valueobject.String1000,
	createdAt time.Time,
	updatedAt time.Time,
) *TrainingSet {
	return &TrainingSet{
		id:                 id,
		trainingExerciseID: trainingExerciseID,
		setNumber:          setNumber,
		weightKg:           weightKg,
		reps:               reps,
		restSeconds:        restSeconds,
		memo:               memo,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
	}
}
