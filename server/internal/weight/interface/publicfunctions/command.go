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

// WeightQuery は weight モジュールが他モジュールに公開する読み取り操作。
type WeightQuery interface {
	CheckIfExistsWeightByUserIDAndMeasuredAt(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error)
}
