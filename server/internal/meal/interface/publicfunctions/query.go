package mealpublicfunctions

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type MealSummaryView struct {
	MealID        valueobject.MealID
	MealType      valueobject.String20
	EatenAt       time.Time
	Calories      valueobject.NonNegativeInt
	ProteinG      *valueobject.NonNegativeDecimal
	FatG          *valueobject.NonNegativeDecimal
	CarbohydrateG *valueobject.NonNegativeDecimal
}

type MealQuery interface {
	ListMealDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error)
	ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*MealSummaryView, error)
}
