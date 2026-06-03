package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type RoutineSpec struct {
	Name      valueobject.String50
	Exercises []RoutineExerciseSpec
}

type RoutineExerciseSpec struct {
	ExerciseID   valueobject.ExerciseID
	DisplayOrder valueobject.NonNegativeInt
}

type Routine struct {
	id        valueobject.RoutineID
	userID    valueobject.UserID
	name      valueobject.String50
	createdAt time.Time
	updatedAt time.Time

	exercises []*RoutineExercise
}

func (r *Routine) ID() valueobject.RoutineID {
	return r.id
}

func (r *Routine) UserID() valueobject.UserID {
	return r.userID
}

func (r *Routine) Name() valueobject.String50 {
	return r.name
}

func (r *Routine) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Routine) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r *Routine) Exercises() []*RoutineExercise {
	return r.exercises
}

func CreateRoutine(spec RoutineSpec, userID valueobject.UserID) *Routine {
	now := time.Now()
	routineID := valueobject.NewPrimaryID[valueobject.RoutineID]()
	return &Routine{
		id:        routineID,
		userID:    userID,
		name:      spec.Name,
		createdAt: now,
		updatedAt: now,
		exercises: rebuildRoutineExercises(routineID, spec.Exercises),
	}
}

func (r *Routine) Update(spec RoutineSpec) {
	r.name = spec.Name
	r.exercises = rebuildRoutineExercises(r.id, spec.Exercises)
	r.updatedAt = time.Now()
}

func NewRoutine(
	id valueobject.RoutineID,
	userID valueobject.UserID,
	name valueobject.String50,
	createdAt time.Time,
	updatedAt time.Time,
	exercises []*RoutineExercise,
) *Routine {
	return &Routine{
		id:        id,
		userID:    userID,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
		exercises: exercises,
	}
}

// spec -> child entities
func rebuildRoutineExercises(routineID valueobject.RoutineID, specs []RoutineExerciseSpec) []*RoutineExercise {
	exerciseRows := make([]*RoutineExercise, 0, len(specs))
	for _, exercise := range specs {
		exerciseRows = append(exerciseRows, CreateRoutineExercise(routineID, exercise.ExerciseID, exercise.DisplayOrder))
	}
	return exerciseRows
}
