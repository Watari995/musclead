package trainingdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type TrainingSummaryView struct {
	TrainingID    valueobject.TrainingID
	StartedAt     time.Time
	EndedAt       *time.Time
	ExerciseCount valueobject.NonNegativeInt // set数が一つ以上の種目が対象(routineから始めるとset数が0でも保存されることがあるため)
	SetCount      valueobject.NonNegativeInt
}

type TrainingQueryService interface {
	ListTrainingDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) (
		[]time.Time, error,
	)
	ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*TrainingSummaryView, error)
	GetTrainingCountInAWeek(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (valueobject.NonNegativeInt, error)
}
