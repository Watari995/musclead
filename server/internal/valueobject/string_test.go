package valueobject_test

import (
	"strings"
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewString20(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "lunch", false},
		{"empty", "", true},
		{"1 char", "a", false},
		{"20 chars", strings.Repeat("a", 20), false},
		{"21 chars", strings.Repeat("a", 21), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewString20(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}

func TestNewString50(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "watari", false},
		{"empty", "", true},
		{"50 chars", strings.Repeat("a", 50), false},
		{"51 chars", strings.Repeat("a", 51), true},
		{"japanese 50 runes", strings.Repeat("あ", 50), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewString50(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}

func TestNewString100(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "hello", false},
		{"empty", "", true},
		{"100 chars", strings.Repeat("a", 100), false},
		{"101 chars", strings.Repeat("a", 101), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewString100(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}

func TestNewString1000(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid memo", "long memo for meal", false},
		{"empty", "", true},
		{"1000 chars", strings.Repeat("a", 1000), false},
		{"1001 chars", strings.Repeat("a", 1001), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewString1000(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}
