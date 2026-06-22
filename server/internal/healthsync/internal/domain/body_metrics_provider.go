package healthsyncdomain

import (
	"context"
	"time"
)

type BodyMetrics struct {
	Weight           float64
	BodyFatPercent   *float64
	SkeletalMuscleKg *float64
	MeasuredAt       time.Time
}

type BodyMetricsProvider interface {
	FetchMetrics(ctx context.Context, accessToken string, from, to time.Time) ([]BodyMetrics, error)
}
