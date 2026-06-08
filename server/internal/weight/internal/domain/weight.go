package weightdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type WeightSpec struct {
	WeightKg          valueobject.WeightKg
	BodyFatPercentage *valueobject.Percentage
	SkeletalMuscleKg  *valueobject.WeightKg
	MeasuredAt        time.Time
}

type Weight struct {
	id                valueobject.WeightID
	userID            valueobject.UserID
	weightKg          valueobject.WeightKg
	bodyFatPercentage *valueobject.Percentage
	skeletalMuscleKg  *valueobject.WeightKg
	measuredAt        time.Time
	createdAt         time.Time
	updatedAt         time.Time
}

func (w *Weight) ID() valueobject.WeightID {
	return w.id
}

func (w *Weight) UserID() valueobject.UserID {
	return w.userID
}

func (w *Weight) WeightKg() valueobject.WeightKg {
	return w.weightKg
}

func (w *Weight) BodyFatPercentage() *valueobject.Percentage {
	return w.bodyFatPercentage
}

func (w *Weight) SkeletalMuscleKg() *valueobject.WeightKg {
	return w.skeletalMuscleKg
}

func (w *Weight) MeasuredAt() time.Time {
	return w.measuredAt
}

func (w *Weight) CreatedAt() time.Time {
	return w.createdAt
}

func (w *Weight) UpdatedAt() time.Time {
	return w.updatedAt
}

func (w *Weight) Update(spec WeightSpec) {
	w.weightKg = spec.WeightKg
	w.bodyFatPercentage = spec.BodyFatPercentage
	w.skeletalMuscleKg = spec.SkeletalMuscleKg
	w.measuredAt = spec.MeasuredAt
	w.updatedAt = time.Now()
}

func CreateWeight(
	userID valueobject.UserID,
	spec WeightSpec,
) *Weight {
	now := time.Now()
	return &Weight{
		id:                valueobject.NewPrimaryID[valueobject.WeightID](),
		userID:            userID,
		weightKg:          spec.WeightKg,
		bodyFatPercentage: spec.BodyFatPercentage,
		skeletalMuscleKg:  spec.SkeletalMuscleKg,
		measuredAt:        spec.MeasuredAt,
		createdAt:         now,
		updatedAt:         now,
	}
}

func NewWeight(
	id valueobject.WeightID,
	userID valueobject.UserID,
	weightKg valueobject.WeightKg,
	bodyFatPercentage *valueobject.Percentage,
	skeletalMuscleKg *valueobject.WeightKg,
	measuredAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Weight {
	return &Weight{
		id:                id,
		userID:            userID,
		weightKg:          weightKg,
		bodyFatPercentage: bodyFatPercentage,
		skeletalMuscleKg:  skeletalMuscleKg,
		measuredAt:        measuredAt,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}
