package weightdto

import (
	"time"

	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type RecordWeightRequest struct {
	WeightKg          string    `json:"weight_kg"`
	BodyFatPercentage *string   `json:"body_fat_percentage,omitempty"`
	SkeletalMuscleKg  *string   `json:"skeletal_muscle_kg,omitempty"`
	MeasuredAt        time.Time `json:"measured_at"`
}

type RecordWeightResponse struct {
	WeightID string `json:"weight_id"`
}

type ListWeightsResponse struct {
	Weights    []WeightDTO             `json:"weights"`
	Pagination shareddto.PaginationDTO `json:"pagination"`
}

type UpdateWeightRequest struct {
	WeightKg          string    `json:"weight_kg"`
	BodyFatPercentage *string   `json:"body_fat_percentage,omitempty"`
	SkeletalMuscleKg  *string   `json:"skeletal_muscle_kg,omitempty"`
	MeasuredAt        time.Time `json:"measured_at"`
}

type UpdateWeightResponse struct {
	WeightID string `json:"weight_id"`
}

type WeightDTO struct {
	ID                string    `json:"id"`
	WeightKg          string    `json:"weight_kg"`
	BodyFatPercentage *string   `json:"body_fat_percentage,omitempty"`
	SkeletalMuscleKg  *string   `json:"skeletal_muscle_kg,omitempty"`
	MeasuredAt        time.Time `json:"measured_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// entity->dtoに変換する
func FromEntity(w *weightdomain.Weight) WeightDTO {
	// nullableなvoをstringに変換する
	var bodyFatPercentage *string
	if w.BodyFatPercentage() != nil {
		s := w.BodyFatPercentage().Value().String()
		bodyFatPercentage = &s
	}
	var skeletalMuscleKg *string
	if w.SkeletalMuscleKg() != nil {
		s := w.SkeletalMuscleKg().Value().String()
		skeletalMuscleKg = &s
	}
	return WeightDTO{
		ID:                w.ID().Value(),
		WeightKg:          w.WeightKg().String(),
		BodyFatPercentage: bodyFatPercentage,
		SkeletalMuscleKg:  skeletalMuscleKg,
		MeasuredAt:        w.MeasuredAt(),
		CreatedAt:         w.CreatedAt(),
		UpdatedAt:         w.UpdatedAt(),
	}
}
