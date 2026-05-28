package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewPercentage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   decimal.Decimal
		wantErr bool
	}{
		{"zero", decimal.Zero, false},
		{"50%", decimal.NewFromFloat(50.5), false},
		{"100%", decimal.NewFromInt(100), false},
		{"negative", decimal.NewFromFloat(-0.01), true},
		{"over 100", decimal.NewFromFloat(100.01), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewPercentage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, tt.input.Equal(got.Value()))
		})
	}
}
