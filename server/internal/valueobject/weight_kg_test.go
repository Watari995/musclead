package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewWeightKg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   decimal.Decimal
		wantErr bool
	}{
		{"valid 70kg", decimal.NewFromFloat(70.5), false},
		{"valid small 0.1kg", decimal.NewFromFloat(0.1), false},
		{"valid 999.99kg", decimal.NewFromFloat(999.99), false},
		{"zero (not allowed)", decimal.Zero, true},
		{"negative", decimal.NewFromFloat(-1), true},
		{"1000kg (excluded)", decimal.NewFromInt(1000), true},
		{"over 1000", decimal.NewFromFloat(1500.5), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewWeightKg(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, tt.input.Equal(got.Value()))
		})
	}
}
