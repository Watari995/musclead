package userdto

import (
	"time"

	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
)

type UpsertWeeklyGoalRequest struct {
	TrainingCount  *int     `json:"training_count"`
	CalorieAverage *int     `json:"calorie_average"`
	WeightChangeKg *float64 `json:"weight_change_kg"`
}

type WeeklyGoalDTO struct {
	TrainingCount  *int      `json:"training_count"`
	CalorieAverage *int      `json:"calorie_average"`
	WeightChangeKg *float64  `json:"weight_change_kg"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func WeeklyGoalFromEntity(g *userdomain.UserWeeklyGoal) *WeeklyGoalDTO {
	if g == nil {
		return nil
	}
	dto := &WeeklyGoalDTO{
		CreatedAt: g.CreatedAt(),
		UpdatedAt: g.UpdatedAt(),
	}
	if g.TrainingCount() != nil {
		v := g.TrainingCount().Value()
		dto.TrainingCount = &v
	}
	if g.CalorieAverage() != nil {
		v := g.CalorieAverage().Value()
		dto.CalorieAverage = &v
	}
	if g.WeightChangeKg() != nil {
		v := g.WeightChangeKg().Value().InexactFloat64()
		dto.WeightChangeKg = &v
	}
	return dto
}
