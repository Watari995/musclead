package valueobject_test

import (
	"testing"

	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"invalid email", "test@example", true},
		{"no at mark", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := valueobject.NewEmail(tt.input)
			// errorを期待していたらassert.Error(t, err)を呼ぶ
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			// errorを期待していなかったらassert.NoError(t, err)を呼ぶ
			assert.NoError(t, err)
			// 結果が期待通りかassert.Equal(t, expected, actual)を呼ぶ
			assert.Equal(t, tt.input, got.Value())
		})
	}

}
