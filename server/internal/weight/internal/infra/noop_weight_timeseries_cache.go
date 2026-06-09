package weightinfra

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type NoOpWeightTimeseriesCache struct{}

func NewNoOpWeightTimeseriesCache() *NoOpWeightTimeseriesCache {
	return &NoOpWeightTimeseriesCache{}
}

func (c *NoOpWeightTimeseriesCache) FindByPeriod(_ context.Context, _ valueobject.UserID, _ time.Time, _ time.Time) ([]*weightdomain.Weight, bool, error) {
	return nil, false, nil
}

func (c *NoOpWeightTimeseriesCache) Add(_ context.Context, _ *weightdomain.Weight) error {
	return nil
}

func (c *NoOpWeightTimeseriesCache) Update(_ context.Context, _ *weightdomain.Weight) error {
	return nil
}

func (c *NoOpWeightTimeseriesCache) Delete(_ context.Context, _ valueobject.UserID, _ valueobject.WeightID) error {
	return nil
}
