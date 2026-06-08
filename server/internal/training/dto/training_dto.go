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

type RecordTrainingSetRequest struct {
	SetNumber   int     `json:"set_number"`
	WeightKg    string  `json:"weight_kg"`
	Reps        int     `json:"reps"`
	RestSeconds *int    `json:"rest_seconds,omitempty"`
	Memo        *string `json:"memo,omitempty"`
}

// req -> spec(vo)に変換する
func (r RecordTrainingSetRequest) ToSpec() (trainingdomain.TrainingSetSpec, error) {
	setNumber, err := valueobject.NewNonNegativeInt(r.SetNumber)
	if err != nil {
		return trainingdomain.TrainingSetSpec{}, myerror.NewBadRequestError().SetMessage("invalid set_number")
	}
	weightKg, err := valueobject.NewNonNegativeDecimalFromString(r.WeightKg)
	if err != nil {
		return trainingdomain.TrainingSetSpec{}, myerror.NewBadRequestError().SetMessage("invalid weight_kg")
	}
	reps, err := valueobject.NewNonNegativeInt(r.Reps)
	if err != nil {
		return trainingdomain.TrainingSetSpec{}, myerror.NewBadRequestError().SetMessage("invalid reps")
	}
	var restSec *valueobject.NonNegativeInt
	if r.RestSeconds != nil {
		r, err := valueobject.NewNonNegativeInt(*r.RestSeconds)
		if err != nil {
			return trainingdomain.TrainingSetSpec{}, myerror.NewBadRequestError().SetMessage("invalid rest_seconds")
		}
		restSec = r
	}
	var memo *valueobject.String1000
	if r.Memo != nil {
		m, err := valueobject.NewString1000(*r.Memo)
		if err != nil {
			return trainingdomain.TrainingSetSpec{}, myerror.NewBadRequestError().SetMessage("invalid memo")
		}
		memo = m
	}
	return trainingdomain.TrainingSetSpec{
		SetNumber:   *setNumber,
		WeightKg:    *weightKg,
		Reps:        *reps,
		RestSeconds: restSec,
		Memo:        memo,
	}, nil
}

type RecordTrainingExerciseRequest struct {
	ExerciseID   string                     `json:"exercise_id"`
	DisplayOrder int                        `json:"display_order"`
	RestSeconds  *int                       `json:"rest_seconds,omitempty"`
	Memo         *string                    `json:"memo,omitempty"`
	Sets         []RecordTrainingSetRequest `json:"sets"`
}

func (r RecordTrainingExerciseRequest) ToSpec() (trainingdomain.TrainingExerciseSpec, error) {
	exerciseID, err := valueobject.NewPrimaryIDFromString[valueobject.ExerciseID](r.ExerciseID)
	if err != nil {
		return trainingdomain.TrainingExerciseSpec{}, myerror.NewBadRequestError().SetMessage("invalid exercise_id")
	}
	displayOrder, err := valueobject.NewNonNegativeInt(r.DisplayOrder)
	if err != nil {
		return trainingdomain.TrainingExerciseSpec{}, myerror.NewBadRequestError().SetMessage("invalid display_order")
	}
	var restSec *valueobject.NonNegativeInt
	if r.RestSeconds != nil {
		r, err := valueobject.NewNonNegativeInt(*r.RestSeconds)
		if err != nil {
			return trainingdomain.TrainingExerciseSpec{}, myerror.NewBadRequestError().SetMessage("invalid rest_seconds")
		}
		restSec = r
	}
	var memo *valueobject.String1000
	if r.Memo != nil {
		m, err := valueobject.NewString1000(*r.Memo)
		if err != nil {
			return trainingdomain.TrainingExerciseSpec{}, myerror.NewBadRequestError().SetMessage("invalid memo")
		}
		memo = m
	}
	sets := lo.Map(r.Sets, func(s RecordTrainingSetRequest, _ int) trainingdomain.TrainingSetSpec {
		spec, err := s.ToSpec()
		if err != nil {
			return trainingdomain.TrainingSetSpec{}
		}
		return spec
	})
	return trainingdomain.TrainingExerciseSpec{
		ExerciseID:   *exerciseID,
		DisplayOrder: *displayOrder,
		RestSeconds:  restSec,
		Memo:         memo,
		Sets:         sets,
	}, nil
}

type RecordTrainingRequest struct {
	StartedAt time.Time               `json:"started_at"`
	EndedAt   *time.Time              `json:"ended_at,omitempty"`
	Memo      *string                 `json:"memo,omitempty"`
	Exercises []RecordTrainingExerciseRequest `json:"exercises"`
}

func (r RecordTrainingRequest) ToSpec() (trainingdomain.TrainingSpec, error) {
	var memo *valueobject.String1000
	if r.Memo != nil {
		m, err := valueobject.NewString1000(*r.Memo)
		if err != nil {
			return trainingdomain.TrainingSpec{}, myerror.NewBadRequestError().SetMessage("invalid memo")
		}
		memo = m
	}
	exercises := lo.Map(r.Exercises, func(e RecordTrainingExerciseRequest, _ int) trainingdomain.TrainingExerciseSpec {
		spec, err := e.ToSpec()
		if err != nil {
			return trainingdomain.TrainingExerciseSpec{}
		}
		return spec
	})
	return trainingdomain.TrainingSpec{
		StartedAt: r.StartedAt,
		EndedAt:   r.EndedAt,
		Memo:      memo,
		Exercises: exercises,
	}, nil
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

func TrainingSetFromEntity(s *trainingdomain.TrainingSet) TrainingSetDTO {
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
	ExerciseID   string           `json:"exercise_id"`
	DisplayOrder int              `json:"display_order"`
	RestSeconds  *int             `json:"rest_seconds,omitempty"`
	Memo         *string          `json:"memo,omitempty"`
	Sets         []TrainingSetDTO `json:"sets"`
}

func TrainingExerciseFromEntity(e *trainingdomain.TrainingExercise) TrainingExerciseDTO {
	// nullableなvoをintに変換
	var restSecondsInt *int
	if e.RestSeconds() != nil {
		r := e.RestSeconds().Value()
		restSecondsInt = &r
	}

	sets := lo.Map(e.Sets(), func(s *trainingdomain.TrainingSet, _ int) TrainingSetDTO {
		return TrainingSetFromEntity(s)
	})

	return TrainingExerciseDTO{
		ID:           e.ID().Value(),
		ExerciseID:   e.ExerciseID().Value(),
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

func TrainingFromEntity(t *trainingdomain.Training) TrainingDTO {
	exercises := lo.Map(t.Exercises(), func(e *trainingdomain.TrainingExercise, _ int) TrainingExerciseDTO {
		return TrainingExerciseFromEntity(e)
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
