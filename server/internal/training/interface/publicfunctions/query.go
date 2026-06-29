package trainingpublicfunctions

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type TrainingSummaryView struct {
	TrainingID    valueobject.TrainingID
	StartedAt     time.Time
	EndedAt       *time.Time
	ExerciseCount valueobject.NonNegativeInt
	SetCount      valueobject.NonNegativeInt
}

type TrainingQuery interface {
	ListTrainingDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error)
	ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*TrainingSummaryView, error)
}
