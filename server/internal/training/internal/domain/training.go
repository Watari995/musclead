package trainingdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

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
	exercises []*TrainingExercise,
) *Training {
	now := time.Now()
	return &Training{
		id:        valueobject.NewPrimaryID[valueobject.TrainingID](),
		userID:    userID,
		startedAt: startedAt,
		endedAt:   endedAt,
		memo:      memo,
		createdAt: now,
		updatedAt: now,
		exercises: exercises,
	}
}