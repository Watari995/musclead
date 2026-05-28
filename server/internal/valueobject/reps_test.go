package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewReps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"1 rep", 1, false},
		{"100 reps", 100, false},
		{"zero (not allowed)", 0, true},
		{"negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewReps(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}
