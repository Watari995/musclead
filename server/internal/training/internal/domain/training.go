package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// ExerciseSpec は新規 Training を組み立てるための「種目1件分の素材」。
// ID / createdAt 等の永続化メタは持たず、 ファクトリ側で付与される。
type ExerciseSpec struct {
	Name         valueobject.String50
	DisplayOrder valueobject.NonNegativeInt
	RestSeconds  *valueobject.NonNegativeInt
	Memo         *valueobject.String1000
	Sets         []SetSpec
}

// SetSpec は ExerciseSpec の中に含めるセット 1 件分の素材。
type SetSpec struct {
	SetNumber   valueobject.NonNegativeInt
	WeightKg    valueobject.NonNegativeDecimal
	Reps        valueobject.NonNegativeInt
	RestSeconds *valueobject.NonNegativeInt
	Memo        *valueobject.String1000
}

type Training struct {
	id        valueobject.TrainingID
	userID    valueobject.UserID
	startedAt time.Time
	endedAt   *time.Time
	memo      *valueobject.String1000
	createdAt time.Time
	updatedAt time.Time

	exercises []*TrainingExercise
}

func (t *Training) ID() valueobject.TrainingID {
	return t.id
}

func (t *Training) UserID() valueobject.UserID {
	return t.userID
}

func (t *Training) StartedAt() time.Time {
	return t.startedAt
}

func (t *Training) EndedAt() *time.Time {
	return t.endedAt
}

func (t *Training) Memo() *valueobject.String1000 {
	return t.memo
}

func (t *Training) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Training) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Training) Exercises() []*TrainingExercise {
	return t.exercises
}

func CreateTraining(
	userID valueobject.UserID,
	startedAt time.Time,
	endedAt *time.Time,
	memo *valueobject.String1000,
	exercises []ExerciseSpec,
) *Training {
	trainingID := valueobject.NewPrimaryID[valueobject.TrainingID]()
	exerciseRows := make([]*TrainingExercise, 0, len(exercises))
	// sets, exerciseを先に作成する
	for _, ex := range exercises {
		exerciseID := valueobject.NewPrimaryID[valueobject.ExerciseID]()
		setRows := make([]*TrainingSet, 0, len(ex.Sets))
		for _, set := range ex.Sets {
			setEntity := CreateTrainingSet(
				exerciseID,
				set.SetNumber,
				set.WeightKg,
				set.Reps,
				set.RestSeconds,
				set.Memo,
			)
			setRows = append(setRows, setEntity)
		}
		exerciseEntity := CreateTrainingExercise(
			trainingID,
			ex.Name,
			ex.DisplayOrder,
			ex.RestSeconds,
			ex.Memo,
			setRows,
		)
		exerciseRows = append(exerciseRows, exerciseEntity)
	}

	now := time.Now()
	return &Training{
		id:        trainingID,
		userID:    userID,
		startedAt: startedAt,
		endedAt:   endedAt,
		memo:      memo,
		createdAt: now,
		updatedAt: now,
		exercises: exerciseRows,
	}
}

func NewTraining(
	id valueobject.TrainingID,
	userID valueobject.UserID,
	startedAt time.Time,
	endedAt *time.Time,
	memo *valueobject.String1000,
	createdAt time.Time,
	updatedAt time.Time,
	exercises []*TrainingExercise,
) *Training {
	return &Training{
		id:        id,
		userID:    userID,
		startedAt: startedAt,
		endedAt:   endedAt,
		memo:      memo,
		createdAt: createdAt,
		updatedAt: updatedAt,
		exercises: exercises,
	}
}
