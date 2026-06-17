package trainingdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type BestSetView struct {
	WeightKg    valueobject.NonNegativeDecimal
	Reps        valueobject.NonNegativeInt
	PerformedAt time.Time
	TrainingID  valueobject.TrainingID
	ExerciseID  valueobject.ExerciseID
}

type ExerciseRecordQueryService interface {
	FindBestSetsByExerciseIDs(ctx context.Context, userID valueobject.UserID, exerciseIDs []valueobject.ExerciseID) ([]*BestSetView, error)
}
