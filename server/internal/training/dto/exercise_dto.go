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

