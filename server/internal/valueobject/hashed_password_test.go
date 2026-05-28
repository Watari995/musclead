package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewHashedPassword(t *testing.T) {
	t.Parallel()

	validHash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid bcrypt hash", string(validHash), false},
		{"empty", "", true},
		{"plain text", "plain-text-password", true},
		{"random string", "abcdef", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewHashedPassword(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}
