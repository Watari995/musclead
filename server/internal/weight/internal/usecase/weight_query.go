package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

type weightQuery struct {
	checkIfExistsWeightByUserIDAndMeasuredAt *CheckIfExistsWeightByUserIDAndMeasuredAt
}

func NewWeightQuery(checkIfExistsWeightByUserIDAndMeasuredAt *CheckIfExistsWeightByUserIDAndMeasuredAt) publicfunctions.WeightQuery {
	return &weightQuery{checkIfExistsWeightByUserIDAndMeasuredAt: checkIfExistsWeightByUserIDAndMeasuredAt}
}

func (q *weightQuery) CheckIfExistsWeightByUserIDAndMeasuredAt(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error) {
	return q.checkIfExistsWeightByUserIDAndMeasuredAt.Execute(ctx, userID, measuredAt)
}
