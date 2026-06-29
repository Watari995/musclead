package publicfunctions

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type WeightRecordInput struct {
	UserID            valueobject.UserID
	WeightKg          valueobject.WeightKg
	BodyFatPercentage *valueobject.Percentage
	SkeletalMuscleKg  *valueobject.WeightKg
	MeasuredAt        time.Time
}

// WeightCommand は weight モジュールが他モジュールに公開する書き込み操作。
type WeightCommand interface {
	Record(ctx context.Context, input WeightRecordInput) (valueobject.WeightID, error)
}

type WeightSummaryView struct {
	WeightID          valueobject.WeightID
	WeightKg          valueobject.WeightKg
	BodyFatPercentage *valueobject.Percentage
	SkeletalMuscleKg  *valueobject.WeightKg
	MeasuredAt        time.Time
}

// WeightQuery は weight モジュールが他モジュールに公開する読み取り操作。
type WeightQuery interface {
	CheckIfExistsWeightByUserIDAndMeasuredAt(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error)
	ListWeightDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error)
	ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*WeightSummaryView, error)
}
