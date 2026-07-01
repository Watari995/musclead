package weightusecase

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type weightQuery struct {
	checkIfExistsWeightByUserIDAndMeasuredAt *CheckIfExistsWeightByUserIDAndMeasuredAt
	listWeightDatesByMonth                   *ListWeightDatesByMonth
	listWeightSummaryByDate                  *ListWeightSummaryByDate
	getWeightChangeInAWeek                   *GetWeightChangeInAWeek
}

func NewWeightQuery(
	checkIfExistsWeightByUserIDAndMeasuredAt *CheckIfExistsWeightByUserIDAndMeasuredAt,
	listWeightDatesByMonth *ListWeightDatesByMonth,
	listWeightSummaryByDate *ListWeightSummaryByDate,
	getWeightChangeInAWeek *GetWeightChangeInAWeek,
) weightpublicfunctions.WeightQuery {
	return &weightQuery{
		checkIfExistsWeightByUserIDAndMeasuredAt: checkIfExistsWeightByUserIDAndMeasuredAt,
		listWeightDatesByMonth:                   listWeightDatesByMonth,
		listWeightSummaryByDate:                  listWeightSummaryByDate,
		getWeightChangeInAWeek:                   getWeightChangeInAWeek,
	}
}

func (q *weightQuery) CheckIfExistsWeightByUserIDAndMeasuredAt(ctx context.Context, userID valueobject.UserID, measuredAt time.Time) (bool, error) {
	return q.checkIfExistsWeightByUserIDAndMeasuredAt.Execute(ctx, userID, measuredAt)
}

func (q *weightQuery) ListWeightDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error) {
	output, err := q.listWeightDatesByMonth.Execute(ctx, ListWeightDatesByMonthInput{
		UserID: userID,
		Year:   year,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}
	return output.Dates, nil
}

func (q *weightQuery) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*weightpublicfunctions.WeightSummaryView, error) {
	output, err := q.listWeightSummaryByDate.Execute(ctx, ListWeightSummaryByDateInput{
		UserID: userID,
		Date:   date,
	})
	if err != nil {
		return nil, err
	}
	return toWeightSummaryViews(output.WeightSummaries), nil
}

func (q *weightQuery) GetWeightChangeInAWeek(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (*valueobject.WeightChangeKg, error) {
	return q.getWeightChangeInAWeek.Execute(ctx, userID, weekStart)
}

func toWeightSummaryViews(views []*weightdomain.WeightSummaryView) []*weightpublicfunctions.WeightSummaryView {
	result := make([]*weightpublicfunctions.WeightSummaryView, 0, len(views))
	for _, v := range views {
		result = append(result, &weightpublicfunctions.WeightSummaryView{
			WeightID:          v.WeightID,
			WeightKg:          v.WeightKg,
			BodyFatPercentage: v.BodyFatPercentage,
			SkeletalMuscleKg:  v.SkeletalMuscleKg,
			MeasuredAt:        v.MeasuredAt,
		})
	}
	return result
}
