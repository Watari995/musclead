package mealusecase

import (
	"context"
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type mealQuery struct {
	listMealDatesByMonth  *ListMealDatesByMonth
	listMealSummaryByDate *ListMealSummaryByDate
}

func NewMealQuery(
	listMealDatesByMonth *ListMealDatesByMonth,
	listMealSummaryByDate *ListMealSummaryByDate,
) mealpublicfunctions.MealQuery {
	return &mealQuery{
		listMealDatesByMonth:  listMealDatesByMonth,
		listMealSummaryByDate: listMealSummaryByDate,
	}
}

func (q *mealQuery) ListMealDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error) {
	output, err := q.listMealDatesByMonth.Execute(ctx, ListMealDatesByMonthInput{
		UserID: userID,
		Year:   year,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}
	return output.Dates, nil
}

func (q *mealQuery) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*mealdomain.MealSummaryView, error) {
	output, err := q.listMealSummaryByDate.Execute(ctx, ListMealSummaryByDateInput{
		UserID: userID,
		Date:   date,
	})
	if err != nil {
		return nil, err
	}
	return output.MealSummaries, nil
}
