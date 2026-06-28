package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type weightQuery struct {
	checkIfExistsWeightByUserIDAndMeasuredAt *CheckIfExistsWeightByUserIDAndMeasuredAt
	listWeightDatesByMonth                   *ListWeightDatesByMonth
	listWeightSummaryByDate                  *ListWeightSummaryByDate
}

func NewWeightQuery(
	checkIfExistsWeightByUserIDAndMeasuredAt *CheckIfExistsWeightByUserIDAndMeasuredAt,
	listWeightDatesByMonth *ListWeightDatesByMonth,
	listWeightSummaryByDate *ListWeightSummaryByDate,
) publicfunctions.WeightQuery {
	return &weightQuery{
		checkIfExistsWeightByUserIDAndMeasuredAt: checkIfExistsWeightByUserIDAndMeasuredAt,
		listWeightDatesByMonth:                   listWeightDatesByMonth,
		listWeightSummaryByDate:                  listWeightSummaryByDate,
	}
}

func (q *weightQuery) CheckIfExistsWeightByUserIDAndMeasuredAt(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error) {
	return q.checkIfExistsWeightByUserIDAndMeasuredAt.Execute(ctx, userID, measuredAt)
}

func (q *weightQuery) ListWeightDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error) {
	input := ListWeightDatesByMonthInput{
		UserID: userID,
		Year:   year,
		Month:  month,
	}
	output, err := q.listWeightDatesByMonth.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.Dates, nil
}

func (q *weightQuery) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*weightdomain.WeightSummaryView, error) {
	input := ListWeightSummaryByDateInput{
		UserID: userID,
		Date:   date,
	}
	output, err := q.listWeightSummaryByDate.Execute(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.WeightSummaries, nil
}
