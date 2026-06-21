package publicfunctions

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type WeightRecordInput struct {
	UserID     valueobject.UserID
	WeightSpec weightdomain.WeightSpec
}

// WeightCommand は weight モジュールが他モジュールに公開する書き込み操作。
type WeightCommand interface {
	Record(ctx context.Context, input WeightRecordInput) (valueobject.WeightID, error)
}

// WeightQuery は weight モジュールが他モジュールに公開する読み取り操作。
type WeightQuery interface {
	CheckIfExistsWeightByUserIDAndMeasuredAt(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error)
}
