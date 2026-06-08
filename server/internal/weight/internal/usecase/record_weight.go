package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
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

func (uc *RecordWeight) Execute(ctx context.Context, input RecordWeightInput) (*RecordWeightOutput, error) {
	weight := weightdomain.CreateWeight(
		input.UserID,
		input.WeightKg,
		input.BodyFatPercentage,
		input.SkeletalMuscleKg,
		input.MeasuredAt,
	)
	if err := uc.weightRepo.Save(ctx, weight); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &RecordWeightOutput{WeightID: weight.ID()}, nil
}

func NewRecordWeight(weightRepo weightdomain.WeightRepository) *RecordWeight {
	return &RecordWeight{weightRepo: weightRepo}
}
