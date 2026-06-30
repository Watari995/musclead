package trainingusecase

import (
	"context"
	"time"

	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type trainingQuery struct {
	listTrainingDatesByMonth  *ListTrainingDatesByMonth
	listTrainingSummaryByDate *ListTrainingSummaryByDate
	countSessionsByWeek       *CountSessionsByWeek
}

func NewTrainingQuery(
	listTrainingDatesByMonth *ListTrainingDatesByMonth,
	listTrainingSummaryByDate *ListTrainingSummaryByDate,
	countSessionsByWeek *CountSessionsByWeek,
) trainingpublicfunctions.TrainingQuery {
	return &trainingQuery{
		listTrainingDatesByMonth:  listTrainingDatesByMonth,
		listTrainingSummaryByDate: listTrainingSummaryByDate,
		countSessionsByWeek:       countSessionsByWeek,
	}
}

func (q *trainingQuery) ListTrainingDatesByMonth(ctx context.Context, userID valueobject.UserID, year, month int) ([]time.Time, error) {
	output, err := q.listTrainingDatesByMonth.Execute(ctx, ListTrainingDatesByMonthInput{
		UserID: userID,
		Year:   year,
		Month:  month,
	})
	if err != nil {
		return nil, err
	}
	return output.Dates, nil
}

func (q *trainingQuery) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*trainingpublicfunctions.TrainingSummaryView, error) {
	output, err := q.listTrainingSummaryByDate.Execute(ctx, ListTrainingSummaryByDateInput{
		UserID: userID,
		Date:   date,
	})
	if err != nil {
		return nil, err
	}
	return toTrainingSummaryViews(output.TrainingSummaries), nil
}

func (q *trainingQuery) CountSessionsByWeek(ctx context.Context, userID valueobject.UserID, weekStart time.Time) (valueobject.NonNegativeInt, error) {
	return q.countSessionsByWeek.Execute(ctx, userID, weekStart)
}

func toTrainingSummaryViews(views []*trainingdomain.TrainingSummaryView) []*trainingpublicfunctions.TrainingSummaryView {
	result := make([]*trainingpublicfunctions.TrainingSummaryView, 0, len(views))
	for _, v := range views {
		result = append(result, &trainingpublicfunctions.TrainingSummaryView{
			TrainingID:    v.TrainingID,
			StartedAt:     v.StartedAt,
			EndedAt:       v.EndedAt,
			ExerciseCount: v.ExerciseCount,
			SetCount:      v.SetCount,
		})
	}
	return result
}
