package trainingdto

import (
	"time"

	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
)

type ExerciseDTO struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ExerciseFromEntity(e *trainingdomain.Exercise) ExerciseDTO {
	return ExerciseDTO{
		ID:           e.ID().Value(),
		UserID:       e.UserID().Value(),
		Name:         e.Name().Value(),
		DisplayOrder: e.DisplayOrder().Value(),
		CreatedAt:    e.CreatedAt(),
		UpdatedAt:    e.UpdatedAt(),
	}
}

type ListExercisesResponse struct {
	Exercises  []ExerciseDTO           `json:"exercises"`
	Pagination shareddto.PaginationDTO `json:"pagination"`
}

type UpsertExerciseRequest struct {
	Name string `json:"name"`
}

type UpsertExerciseResponse struct {
	ID string `json:"id"`
}

type ReorderExercisesRequest struct {
	ExerciseIDs []string `json:"exercise_ids"`
}

type BestSetDTO struct {
	WeightKg    string    `json:"weight_kg"` // weightは精度のためstringとして持つ
	Reps        int       `json:"reps"`
	PerformedAt time.Time `json:"performed_at"`
	TrainingID  string    `json:"training_id"`
	ExerciseID  string    `json:"exercise_id"`
}

func BestSetFromData(b *trainingdomain.BestSetView) BestSetDTO {
	return BestSetDTO{
		WeightKg:    b.WeightKg.String(),
		Reps:        b.Reps.Value(),
		PerformedAt: b.PerformedAt,
		TrainingID:  b.TrainingID.Value(),
		ExerciseID:  b.ExerciseID.Value(),
	}
}

type ListBestSetsResponse struct {
	BestSets []BestSetDTO `json:"best_sets"`
}

// BestSetTimeseriesDataPointDTO は1セッション分のベストセットを表す時系列の1点。
// BestSetDTO と同じフィールドだが、timeseries 文脈で名前を明確にしている。
type BestSetTimeseriesDataPointDTO struct {
	PerformedAt time.Time `json:"performed_at"`
	WeightKg    string    `json:"weight_kg"`
	Reps        int       `json:"reps"`
	TrainingID  string    `json:"training_id"`
}

func BestSetTimeseriesDataPointFromData(b *trainingdomain.BestSetView) BestSetTimeseriesDataPointDTO {
	return BestSetTimeseriesDataPointDTO{
		PerformedAt: b.PerformedAt,
		WeightKg:    b.WeightKg.String(),
		Reps:        b.Reps.Value(),
		TrainingID:  b.TrainingID.Value(),
	}
}

type BestSetTimeseriesResponse struct {
	Period     string                           `json:"period"`
	ExerciseID string                           `json:"exercise_id"`
	DataPoints []BestSetTimeseriesDataPointDTO  `json:"data_points"`
}
