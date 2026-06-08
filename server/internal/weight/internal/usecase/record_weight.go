package weightusecase

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type RecordWeightInput struct {
	UserID            valueobject.UserID
	WeightKg          valueobject.WeightKg
	BodyFatPercentage *valueobject.Percentage
	SkeletalMuscleKg  *valueobject.WeightKg
	MeasuredAt        time.Time
}

type RecordWeightOutput struct {
	WeightID valueobject.WeightID
}

type RecordWeight struct {
	weightRepo weightdomain.WeightRepository
}
