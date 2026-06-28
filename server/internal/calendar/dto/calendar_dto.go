package calendardto

import "time"

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

type MealSummaryDTO struct {
	MealID        string  `json:"meal_id"`
	MealType      string  `json:"meal_type"`
	EatenAt       string  `json:"eaten_at"`
	Calories      int     `json:"calories"`
	ProteinG      *string `json:"protein_g,omitempty"`
	FatG          *string `json:"fat_g,omitempty"`
	CarbohydrateG *string `json:"carbohydrate_g,omitempty"`
}

type WeightSummaryDTO struct {
	WeightID          string  `json:"weight_id"`
	WeightKg          string  `json:"weight_kg"`
	BodyFatPercentage *string `json:"body_fat_percentage,omitempty"`
	SkeletalMuscleKg  *string `json:"skeletal_muscle_kg,omitempty"`
	MeasuredAt        string  `json:"measured_at"`
}

type GetDailySummaryResponse struct {
	Trainings []TrainingSummaryDTO `json:"trainings"`
	Meals     []MealSummaryDTO     `json:"meals"`
	Weights   []WeightSummaryDTO   `json:"weights"`
}
