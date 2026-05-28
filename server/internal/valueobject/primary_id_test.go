package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPrimaryId_GeneratesUUIDv7(t *testing.T) {
	t.Parallel()

	id := valueobject.NewPrimaryId[valueobject.UserID]()

	// UUID v7 として parse できる
	parsed, err := uuid.Parse(id.Value())
	assert.NoError(t, err)
	assert.Equal(t, uuid.Version(7), parsed.Version())
}

func TestNewPrimaryId_UniquePerCall(t *testing.T) {
	t.Parallel()

	a := valueobject.NewPrimaryId[valueobject.UserID]()
	b := valueobject.NewPrimaryId[valueobject.UserID]()

	assert.NotEqual(t, a.Value(), b.Value())
}

func TestNewPrimaryIdFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid uuid", "019e6ce9-5253-7816-8af8-6603ae0516e6", false},
		{"empty", "", true},
		{"not a uuid", "not-a-uuid", true},
		{"missing hyphens", "019e6ce952537816", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewPrimaryIdFromString[valueobject.UserID](tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.input, got.Value())
		})
	}
}

func TestPrimaryId_Bytes(t *testing.T) {
	t.Parallel()

	id, err := valueobject.NewPrimaryIdFromString[valueobject.UserID]("019e6ce9-5253-7816-8af8-6603ae0516e6")
	assert.NoError(t, err)

	b, err := id.Bytes()
	assert.NoError(t, err)
	assert.Len(t, b, 16) // UUID は 16 バイト
}
