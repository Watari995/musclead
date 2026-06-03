package trainingdto

import (
	"time"

	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/samber/lo"
)

// ─── Request / Response (HTTP境界) ─────────────────────────

type ListRoutinesResponse struct {
	Routines   []RoutineDTO            `json:"routines"`
	Pagination shareddto.PaginationDTO `json:"pagination"`
}

type UpsertRoutineExerciseRequest struct {
	ExerciseID   string `json:"exercise_id"`
	DisplayOrder int    `json:"display_order"`
}
type UpsertRoutineRequest struct {
	Name      string                         `json:"name"`
	Exercises []UpsertRoutineExerciseRequest `json:"exercises"`
}
type UpsertRoutineResponse struct {
	ID string `json:"id"`
}

// ----- Entity view ────────────────────────────────────────

type RoutineExerciseDTO struct {
	ID           string `json:"id"`
	ExerciseID   string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
	DisplayOrder int    `json:"display_order"`
}

func NewRoutineExerciseDTO(e *trainingdomain.RoutineExerciseView) RoutineExerciseDTO {
	return RoutineExerciseDTO{
		ID:           e.ID.Value(),
		ExerciseID:   e.ExerciseID.Value(),
		ExerciseName: e.ExerciseName.Value(),
		DisplayOrder: e.DisplayOrder.Value(),
	}
}

type RoutineDTO struct {
	ID               string               `json:"id"`
	UserID           string               `json:"user_id"`
	Name             string               `json:"name"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	RoutineExercises []RoutineExerciseDTO `json:"routine_exercises"`
}

func NewRoutineDTO(r *trainingdomain.RoutineView) RoutineDTO {
	return RoutineDTO{
		ID:        r.ID.Value(),
		UserID:    r.UserID.Value(),
		Name:      r.Name.Value(),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		RoutineExercises: lo.Map(r.RoutineExercises, func(e trainingdomain.RoutineExerciseView, _ int) RoutineExerciseDTO {
			return NewRoutineExerciseDTO(&e)
		}),
	}
}
