package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// ExerciseSpec は新規 Training を組み立てるための「種目1件分の素材」。
// ID / createdAt 等の永続化メタは持たず、 ファクトリ側で付与される。
// DDD的にはDataよりもSpecの方がいいらしい。

type TrainingSpec struct {
	StartedAt time.Time
	EndedAt   *time.Time
	Memo      *valueobject.String1000
	Exercises []ExerciseSpec
}
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

func CreateTraining(spec TrainingSpec, userID valueobject.UserID) *Training {
	trainingID := valueobject.NewPrimaryID[valueobject.TrainingID]()
	exerciseRows := rebuildExercises(trainingID, spec.Exercises)
	now := time.Now()
	return &Training{
		id:        trainingID,
		userID:    userID,
		startedAt: spec.StartedAt,
		endedAt:   spec.EndedAt,
		memo:      spec.Memo,
		createdAt: now,
		updatedAt: now,
		exercises: exerciseRows,
	}
}

type UpdateParams struct {
	StartedAt time.Time
	EndedAt   *time.Time
	Memo      *valueobject.String1000
	Exercises []ExerciseSpec
}

func (t *Training) Update(params UpdateParams) {
	t.startedAt = params.StartedAt
	t.endedAt = params.EndedAt
	t.memo = params.Memo
	t.exercises = rebuildExercises(t.id, params.Exercises)
	t.updatedAt = time.Now()
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

func rebuildExercises(trainingID valueobject.TrainingID, specs []ExerciseSpec) []*TrainingExercise {
	// sets, exerciseを先に作成する
	exerciseRows := make([]*TrainingExercise, 0, len(specs))
	for _, ex := range specs {
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
	return exerciseRows
}
