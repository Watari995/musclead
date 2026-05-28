package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewNonNegativeInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"zero", 0, false},
		{"positive", 100, false},
		{"large positive", 1_000_000, false},
		{"negative", -1, true},
		{"very negative", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewNonNegativeInt(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}
