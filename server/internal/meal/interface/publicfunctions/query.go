package mealpublicfunctions

import (
	"context"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type MealQuery interface {
	ListMealDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error)
	ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*mealdomain.MealSummaryView, error)
}
