package shareddto_test

import (
	"encoding/json"
	"testing"

	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/stretchr/testify/assert"
)

type wrap struct {
	X shareddto.Patch[string] `json:"x"`
}

func TestPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantSet   bool
		wantNull  bool
		wantValue string
		wantErr   bool
	}{
		{"key absent", `{}`, false, false, "", false},
		{"null", `{"x": null}`, true, true, "", false},
		{"value", `{"x": "abc"}`, true, false, "abc", false},
		{"type mismatch", `{"x": 123}`, false, false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var w wrap
			err := json.Unmarshal([]byte(tt.input), &w)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSet, w.X.Set)
			assert.Equal(t, tt.wantNull, w.X.Null)
			assert.Equal(t, tt.wantValue, w.X.Value)
		})
	}
}
