package trainingusecase

import (
	"context"
	"time"

	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type trainingQuery struct {
	listTrainingDatesByMonth   *ListTrainingDatesByMonth
	listTrainingSummaryByDate  *ListTrainingSummaryByDate
}

func NewTrainingQuery(
	listTrainingDatesByMonth *ListTrainingDatesByMonth,
	listTrainingSummaryByDate *ListTrainingSummaryByDate,
) trainingpublicfunctions.TrainingQuery {
	return &trainingQuery{
		listTrainingDatesByMonth:  listTrainingDatesByMonth,
		listTrainingSummaryByDate: listTrainingSummaryByDate,
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

func (q *trainingQuery) ListSummaryByDate(ctx context.Context, userID valueobject.UserID, date time.Time) ([]*trainingdomain.TrainingSummaryView, error) {
	output, err := q.listTrainingSummaryByDate.Execute(ctx, ListTrainingSummaryByDateInput{
		UserID: userID,
		Date:   date,
	})
	if err != nil {
		return nil, err
	}
	return output.TrainingSummaries, nil
}
