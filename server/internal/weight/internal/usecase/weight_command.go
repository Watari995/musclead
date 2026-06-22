package weightusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type weightCommand struct {
	record *RecordWeight
}

func NewWeightCommand(record *RecordWeight) weightpublicfunctions.WeightCommand {
	return &weightCommand{record: record}
}

func (c *weightCommand) Record(ctx context.Context, input weightpublicfunctions.WeightRecordInput) (valueobject.WeightID, error) {
	output, err := c.record.Execute(ctx, RecordWeightInput{
		UserID: input.UserID,
		WeightSpec: weightdomain.WeightSpec{
			WeightKg:          input.WeightKg,
			BodyFatPercentage: input.BodyFatPercentage,
			SkeletalMuscleKg:  input.SkeletalMuscleKg,
			MeasuredAt:        input.MeasuredAt,
		},
	})
	if err != nil {
		return valueobject.WeightID{}, err
	}
	return output.WeightID, nil
}
