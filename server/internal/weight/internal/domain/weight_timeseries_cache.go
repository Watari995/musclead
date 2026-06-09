package weightdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type WeightTimeseriesCache interface {
	FindByPeriod(ctx context.Context, userID valueobject.UserID, from, to time.Time) (weights []*Weight, hit bool, err error)
	Save(ctx context.Context, weight *Weight) error
	Delete(ctx context.Context, userID valueobject.UserID, weightID valueobject.WeightID) error
}
