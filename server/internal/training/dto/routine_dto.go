package trainingdto

import (
	"time"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
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

func (r UpsertRoutineRequest) ToSpec() ([]trainingdomain.RoutineExerciseSpec, error) {
	specs := make([]trainingdomain.RoutineExerciseSpec, 0, len(r.Exercises))
	for _, e := range r.Exercises {
		exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](e.ExerciseID)
		if err != nil {
			return []trainingdomain.RoutineExerciseSpec{}, myerror.NewBadRequestError().SetMessage("invalid exerciseID")
		}
		displayOrder, err := valueobject.NewNonNegativeInt(e.DisplayOrder)
		if err != nil {
			return []trainingdomain.RoutineExerciseSpec{}, myerror.NewBadRequestError().SetMessage("invalid display order")
		}
		specs = append(specs, trainingdomain.RoutineExerciseSpec{ExerciseID: *exerciseID, DisplayOrder: *displayOrder})
	}
	return specs, nil
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

func RoutineExerciseFromEntity(e *trainingdomain.RoutineExerciseView) RoutineExerciseDTO {
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

func RoutineFromEntity(r *trainingdomain.RoutineView) RoutineDTO {
	return RoutineDTO{
		ID:        r.ID.Value(),
		UserID:    r.UserID.Value(),
		Name:      r.Name.Value(),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		RoutineExercises: lo.Map(r.RoutineExercises, func(e trainingdomain.RoutineExerciseView, _ int) RoutineExerciseDTO {
			return RoutineExerciseFromEntity(&e)
		}),
	}
}
