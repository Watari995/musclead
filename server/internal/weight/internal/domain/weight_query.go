package weightdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type WeightSummaryView struct {
	WeightID          valueobject.WeightID
	WeightKg          valueobject.WeightKg
	BodyFatPercentage *valueobject.Percentage
	SkeletalMuscleKg  *valueobject.WeightKg
	MeasuredAt        time.Time
}

type WeightQueryService interface {
	ListWeightDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) (
		[]time.Time, error,
	)
	ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*WeightSummaryView, error)
}
