package weightdto

import (
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/Watari995/musclead/internal/valueobject"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
)

type UpsertWeightRequest struct {
	WeightKg          string    `json:"weight_kg"`
	BodyFatPercentage *string   `json:"body_fat_percentage,omitempty"`
	SkeletalMuscleKg  *string   `json:"skeletal_muscle_kg,omitempty"`
	MeasuredAt        time.Time `json:"measured_at"`
}

func (r UpsertWeightRequest) ToSpec() (weightdomain.WeightSpec, error) {
	weightKg, err := valueobject.NewWeightKgFromString(r.WeightKg)
	if err != nil {
		return weightdomain.WeightSpec{}, myerror.NewBadRequestError().SetMessage("invalid weight kg")
	}
	var bodyFatPercentage *valueobject.Percentage
	if r.BodyFatPercentage != nil {
		bodyFatPercentage, err = valueobject.NewPercentageFromString(*r.BodyFatPercentage)
		if err != nil {
			return weightdomain.WeightSpec{}, myerror.NewBadRequestError().SetMessage("invalid body fat percentage")
		}
	}
	var skeletalMuscleKg *valueobject.WeightKg
	if r.SkeletalMuscleKg != nil {
		skeletalMuscleKg, err = valueobject.NewWeightKgFromString(*r.SkeletalMuscleKg)
		if err != nil {
			return weightdomain.WeightSpec{}, myerror.NewBadRequestError().SetMessage("invalid skeletal muscle kg")
		}
	}
	return weightdomain.WeightSpec{
		WeightKg:          *weightKg,
		BodyFatPercentage: bodyFatPercentage,
		SkeletalMuscleKg:  skeletalMuscleKg,
		MeasuredAt:        r.MeasuredAt,
	}, nil
}

type UpsertWeightResponse struct {
	WeightID string `json:"weight_id"`
}

type ListWeightsResponse struct {
	Weights    []WeightDTO             `json:"weights"`
	Pagination shareddto.PaginationDTO `json:"pagination"`
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
