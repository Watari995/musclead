package calendardto

import (
	"time"

	mealpublicfunctions "github.com/Watari995/musclead/internal/meal/interface/publicfunctions"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
)

type MonthlySummaryDayDTO struct {
	Date        string `json:"date"`
	HasTraining bool   `json:"has_training"`
	HasMeal     bool   `json:"has_meal"`
	HasWeight   bool   `json:"has_weight"`
}

type GetMonthlySummaryResponse struct {
	Days []MonthlySummaryDayDTO `json:"days"`
}

type TrainingSummaryDTO struct {
	TrainingID    string     `json:"training_id"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
	ExerciseCount int        `json:"exercise_count"`
	SetCount      int        `json:"set_count"`
}

func FromTrainingSummaryViewToDTO(view *trainingpublicfunctions.TrainingSummaryView) TrainingSummaryDTO {
	return TrainingSummaryDTO{
		TrainingID:    view.TrainingID.String(),
		StartedAt:     view.StartedAt,
		EndedAt:       view.EndedAt,
		ExerciseCount: view.ExerciseCount.Value(),
		SetCount:      view.SetCount.Value(),
	}
}

type MealSummaryDTO struct {
	MealID        string  `json:"meal_id"`
	MealType      string  `json:"meal_type"`
	EatenAt       string  `json:"eaten_at"`
	Calories      int     `json:"calories"`
	ProteinG      *string `json:"protein_g,omitempty"`
	FatG          *string `json:"fat_g,omitempty"`
	CarbohydrateG *string `json:"carbohydrate_g,omitempty"`
}

func FromMealSummaryViewToDTO(view *mealpublicfunctions.MealSummaryView) MealSummaryDTO {
	var proteinG, fatG, carbohydrateG *string
	if view.ProteinG != nil {
		proteinStr := view.ProteinG.String()
		proteinG = &proteinStr
	}
	if view.FatG != nil {
		fatStr := view.FatG.String()
		fatG = &fatStr
	}
	if view.CarbohydrateG != nil {
		carbohydrateStr := view.CarbohydrateG.String()
		carbohydrateG = &carbohydrateStr
	}

	return MealSummaryDTO{
		MealID:        view.MealID.String(),
		MealType:      view.MealType.String(),
		EatenAt:       view.EatenAt.Format(time.RFC3339),
		Calories:      view.Calories.Value(),
		ProteinG:      proteinG,
		FatG:          fatG,
		CarbohydrateG: carbohydrateG,
	}
}

type WeightSummaryDTO struct {
	WeightID          string  `json:"weight_id"`
	WeightKg          string  `json:"weight_kg"`
	BodyFatPercentage *string `json:"body_fat_percentage,omitempty"`
	SkeletalMuscleKg  *string `json:"skeletal_muscle_kg,omitempty"`
	MeasuredAt        string  `json:"measured_at"`
}

func FromWeightSummaryViewToDTO(view *publicfunctions.WeightSummaryView) WeightSummaryDTO {
	var bodyFatPercentage, skeletalMuscleKg *string
	if view.BodyFatPercentage != nil {
		bodyFatStr := view.BodyFatPercentage.String()
		bodyFatPercentage = &bodyFatStr
	}
	if view.SkeletalMuscleKg != nil {
		skeletalMuscleStr := view.SkeletalMuscleKg.String()
		skeletalMuscleKg = &skeletalMuscleStr
	}

	return WeightSummaryDTO{
		WeightID:          view.WeightID.String(),
		WeightKg:          view.WeightKg.String(),
		BodyFatPercentage: bodyFatPercentage,
		SkeletalMuscleKg:  skeletalMuscleKg,
		MeasuredAt:        view.MeasuredAt.Format(time.RFC3339),
	}
}

type GetDailySummaryResponse struct {
	Trainings     []TrainingSummaryDTO `json:"trainings"`
	Meals         []MealSummaryDTO     `json:"meals"`
	TotalCalories int                  `json:"total_calories"`
	Weights       []WeightSummaryDTO   `json:"weights"`
}
