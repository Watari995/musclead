package valueobject

import (
	"errors"
	"time"
)

type Period string

const (
	Period1Week    Period = "1week"
	Period1Month   Period = "1month"
	Period3Months  Period = "3months"
	PeriodHalfYear Period = "halfyear"
	Period1Year    Period = "1year"
)

var ErrInvalidPeriod = errors.New("invalid period")

func NewPeriodFromString(s string) (Period, error) {
	switch Period(s) {
	case Period1Week, Period1Month, Period3Months, PeriodHalfYear, Period1Year:
		return Period(s), nil
	default:
		return "", ErrInvalidPeriod
	}
}

func (p Period) Duration() time.Duration {
	switch p {
	case Period1Week:
		return 7 * 24 * time.Hour
	case Period1Month:
		return 30 * 24 * time.Hour
	case Period3Months:
		return 90 * 24 * time.Hour
	case PeriodHalfYear:
		return 180 * 24 * time.Hour
	case Period1Year:
		return 365 * 24 * time.Hour
	default:
		return 0
	}
}

func (p Period) String() string {
	return string(p)
}
