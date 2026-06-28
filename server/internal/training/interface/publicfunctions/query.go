package trainingpublicfunctions

import (
	"context"
	"time"

	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type TrainingQuery interface {
	ListTrainingDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error)
	ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*trainingdomain.TrainingSummaryView, error)
}
