package valueobject

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Reps struct {
	LiteralBase[int]
}

var repsValidationRules = []validation.Rule{
	validation.Required,
	validation.Min(1),
}

func (r Reps) Validate() error {
	return validation.Validate(r.Value(), repsValidationRules...)
}

func NewReps(v int) (*Reps, error) {
	r := Reps{LiteralBase: LiteralBase[int]{v: v}}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &r, nil
}
