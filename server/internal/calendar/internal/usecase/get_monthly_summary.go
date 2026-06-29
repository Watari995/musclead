package calendarusecase

import (
	"context"
	"sort"
	"sync"
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
	weightpublicfunctions "github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

type GetMonthlySummaryInput struct {
	UserID valueobject.UserID
	Year   int
	Month  int
}

type DaySummary struct {
	Date        time.Time
	HasTraining bool
	HasMeal     bool
	HasWeight   bool
}

type GetMonthlySummaryOutput struct {
	Days []DaySummary
}

type GetMonthlySummary struct {
	trainingQuery trainingpublicfunctions.TrainingQuery
	mealQuery     mealpublicfunctions.MealQuery
	weightQuery   weightpublicfunctions.WeightQuery
}

func NewGetMonthlySummary(
	trainingQuery trainingpublicfunctions.TrainingQuery,
	mealQuery mealpublicfunctions.MealQuery,
	weightQuery weightpublicfunctions.WeightQuery,
) *GetMonthlySummary {
	return &GetMonthlySummary{
		trainingQuery: trainingQuery,
		mealQuery:     mealQuery,
		weightQuery:   weightQuery,
	}
}

func (uc *GetMonthlySummary) Execute(ctx context.Context, input GetMonthlySummaryInput) (*GetMonthlySummaryOutput, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var trainingDates []time.Time
	var trainingErr error
	go func() {
		defer wg.Done()
		trainingDates, trainingErr = uc.trainingQuery.ListTrainingDatesByMonth(ctx, input.UserID, input.Year, input.Month)
	}()

	wg.Add(1)
	var mealDates []time.Time
	var mealErr error
	go func() {
		defer wg.Done()
		mealDates, mealErr = uc.mealQuery.ListMealDatesByMonth(ctx, input.UserID, input.Year, input.Month)
	}()

	wg.Add(1)
	var weightDates []time.Time
	var weightErr error
	go func() {
		defer wg.Done()
		weightDates, weightErr = uc.weightQuery.ListWeightDatesByMonth(ctx, input.UserID, input.Year, input.Month)
	}()
	// 全てのgoroutineが終了するまで後続のmapの処理を待機する
	wg.Wait()

	// mapを作成して日付がkey,DaySummaryがvalueのmapを作成
	daySummaryMap := make(map[string]*DaySummary)
	for _, t := range trainingDates {
		getOrCreateDaySummary(t, daySummaryMap).HasTraining = true
	}
	for _, m := range mealDates {
		getOrCreateDaySummary(m, daySummaryMap).HasMeal = true
	}
	for _, w := range weightDates {
		getOrCreateDaySummary(w, daySummaryMap).HasWeight = true
	}

	days := make([]DaySummary, 0, len(daySummaryMap))
	for _, daySummary := range daySummaryMap {
		days = append(days, *daySummary)
	}

	// 日付を昇順にするために、日付をキーにしてソートする
	sort.Slice(days, func(i, j int) bool {
		return days[i].Date.Before(days[j].Date)
	})

	output := &GetMonthlySummaryOutput{
		Days: days,
	}
	if trainingErr != nil {
		return output, trainingErr
	}
	if mealErr != nil {
		return output, mealErr
	}
	if weightErr != nil {
		return output, weightErr
	}

	return output, nil
}

// keyが存在しなければ作成して返却する
func getOrCreateDaySummary(date time.Time, daySummaryMap map[string]*DaySummary) *DaySummary {
	dateStr := date.Format("2006-01-02")
	if _, exists := daySummaryMap[dateStr]; !exists {
		daySummaryMap[dateStr] = &DaySummary{
			Date: date,
		}
	}
	return daySummaryMap[dateStr]
}
