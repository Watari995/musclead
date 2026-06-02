package trainingdto

import (
	"time"

	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
)

// ─── Request / Response (HTTP境界) ─────────────────────────

type RecordSetRequest struct {
	SetNumber   int     `json:"set_number"`
	WeightKg    string  `json:"weight_kg"`
	Reps        int     `json:"reps"`
	RestSeconds *int    `json:"rest_seconds,omitempty"`
	Memo        *string `json:"memo,omitempty"`
}

type RecordExerciseRequest struct {
	Name         string             `json:"name"`
	DisplayOrder int                `json:"display_order"`
	RestSeconds  *int               `json:"rest_seconds,omitempty"`
	Memo         *string            `json:"memo,omitempty"`
	Sets         []RecordSetRequest `json:"sets"`
}

type RecordTrainingRequest struct {
	StartedAt time.Time               `json:"started_at"`
	EndedAt   *time.Time              `json:"ended_at,omitempty"`
	Memo      *string                 `json:"memo,omitempty"`
	Exercises []RecordExerciseRequest `json:"exercises"`
}

type RecordTrainingResponse struct {
	TrainingID string `json:"training_id"`
}

// UpdateTrainingRequest は Record と同じ shape のため alias 的に使い回す。
type UpdateTrainingRequest = RecordTrainingRequest

type UpdateTrainingResponse struct {
	TrainingID string `json:"training_id"`
}

type ListTrainingsResponse struct {
	Trainings  []TrainingDTO           `json:"trainings"`
	Pagination shareddto.PaginationDTO `json:"pagination"`
}

// ─── Entity view ────────────────────────────────────────

type TrainingSetDTO struct {
	ID          string  `json:"id"`
	SetNumber   int     `json:"set_number"`
	WeightKg    string  `json:"weight_kg"` // weightは精度のためstringとして持つ
	Reps        int     `json:"reps"`
	RestSeconds *int    `json:"rest_seconds,omitempty"`
	Memo        *string `json:"memo,omitempty"`
}

func NewTrainingSetDTO(s *trainingdomain.TrainingSet) TrainingSetDTO {
	// nullableなvoをintに変換
	var restSecondsInt *int
	if s.RestSeconds() != nil {
		r := s.RestSeconds().Value()
		restSecondsInt = &r
	}
	return TrainingSetDTO{
		ID:          s.ID().Value(),
		SetNumber:   s.SetNumber().Value(),
		WeightKg:    s.WeightKg().String(),
		Reps:        s.Reps().Value(),
		RestSeconds: restSecondsInt,
		Memo:        memoToPtrStr(s.Memo()),
	}
}

type TrainingExerciseDTO struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	DisplayOrder int              `json:"display_order"`
	RestSeconds  *int             `json:"rest_seconds,omitempty"`
	Memo         *string          `json:"memo,omitempty"`
	Sets         []TrainingSetDTO `json:"sets"`
}

func NewTrainingExerciseDTO(e *trainingdomain.TrainingExercise) TrainingExerciseDTO {
	// nullableなvoをintに変換
	var restSecondsInt *int
	if e.RestSeconds() != nil {
		r := e.RestSeconds().Value()
		restSecondsInt = &r
	}

	sets := lo.Map(e.Sets(), func(s *trainingdomain.TrainingSet, _ int) TrainingSetDTO {
		return NewTrainingSetDTO(s)
	})

	return TrainingExerciseDTO{
		ID:           e.ID().Value(),
		Name:         e.Name().Value(),
		DisplayOrder: e.DisplayOrder().Value(),
		RestSeconds:  restSecondsInt,
		Memo:         memoToPtrStr(e.Memo()),
		Sets:         sets,
	}
}

type TrainingDTO struct {
	ID        string                `json:"id"`
	UserID    string                `json:"user_id"`
	StartedAt time.Time             `json:"started_at"`
	EndedAt   *time.Time            `json:"ended_at"`
	Memo      *string               `json:"memo,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Exercises []TrainingExerciseDTO `json:"exercises"`
}

func NewTrainingDTO(t *trainingdomain.Training) TrainingDTO {
	exercises := lo.Map(t.Exercises(), func(e *trainingdomain.TrainingExercise, _ int) TrainingExerciseDTO {
		return NewTrainingExerciseDTO(e)
	})

	return TrainingDTO{
		ID:        t.ID().Value(),
		UserID:    t.UserID().Value(),
		StartedAt: t.StartedAt(),
		EndedAt:   t.EndedAt(),
		Memo:      memoToPtrStr(t.Memo()),
		CreatedAt: t.CreatedAt(),
		UpdatedAt: t.UpdatedAt(),
		Exercises: exercises,
	}
}

// private 変換汎用メソッド
func memoToPtrStr(memo *valueobject.String1000) *string {
	// nullableなvoをstringに変換
	var memoStr *string
	if memo != nil {
		m := memo.Value()
		memoStr = &m
	}
	return memoStr
}
