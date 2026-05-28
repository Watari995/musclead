package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewNonNegativeDecimal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   decimal.Decimal
		wantErr bool
	}{
		{"zero", decimal.Zero, false},
		{"positive int", decimal.NewFromInt(100), false},
		{"positive float", decimal.NewFromFloat(12.5), false},
		{"negative", decimal.NewFromInt(-1), true},
		{"negative float", decimal.NewFromFloat(-0.01), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewNonNegativeDecimal(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.True(t, tt.input.Equal(got.Value()))
		})
	}
}

func TestParseOptionalNonNegativeDecimal(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		got, err := valueobject.ParseOptionalNonNegativeDecimal(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("valid float", func(t *testing.T) {
		t.Parallel()
		v := 12.5
		got, err := valueobject.ParseOptionalNonNegativeDecimal(&v)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.True(t, decimal.NewFromFloat(12.5).Equal(got.Value()))
	})

	t.Run("negative float returns error", func(t *testing.T) {
		t.Parallel()
		v := -1.0
		got, err := valueobject.ParseOptionalNonNegativeDecimal(&v)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
