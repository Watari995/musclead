package valueobject_test

import (
	"strings"
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"https URL", "https://example.com/path?q=1", false},
		{"http URL", "http://example.com", false},
		{"empty", "", true},
		{"not a URL", "not-a-url", true},
		{"too long", "https://example.com/" + strings.Repeat("a", 2048), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewURL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}
