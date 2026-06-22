package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type CheckIfExistsWeightByUserIDAndMeasuredAt struct {
	weightRepo weightdomain.WeightRepository
}

func (uc *CheckIfExistsWeightByUserIDAndMeasuredAt) Execute(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error) {
	exists, err := uc.weightRepo.ExistsByUserIDAndMeasuredAt(ctx, userID, measuredAt)
	if err != nil {
		return false, myerror.NewInternalError().Wrap(err)
	}
	return exists, nil
}

func NewCheckIfExistsWeightByUserIDAndMeasuredAt(weightRepo weightdomain.WeightRepository) *CheckIfExistsWeightByUserIDAndMeasuredAt {
	return &CheckIfExistsWeightByUserIDAndMeasuredAt{
		weightRepo: weightRepo,
	}
}
