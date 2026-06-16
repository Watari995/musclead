package trainingdto_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	trainingdto "github.com/Watari995/musclead/internal/training/dto"
)

// 不正なセット(weight_kg がパース不能)を含む場合、 ToSpec はエラーを返す。
// 以前は lo.Map がエラーを握り潰し、 0kg×0reps のゼロ値セットを黙って永続化していた。
func TestRecordTrainingRequest_ToSpec_InvalidSetReturnsError(t *testing.T) {
	t.Parallel()

	req := trainingdto.RecordTrainingRequest{
		Exercises: []trainingdto.RecordTrainingExerciseRequest{
			{
				ExerciseID:   uuid.NewString(),
				DisplayOrder: 1,
				Sets: []trainingdto.RecordTrainingSetRequest{
					{SetNumber: 1, WeightKg: "not-a-number", Reps: 12},
				},
			},
		},
	}

	_, err := req.ToSpec()
	assert.Error(t, err)
}

// 有効な入力(9kg×12reps)は値を失わずに Spec へ変換される。
func TestRecordTrainingRequest_ToSpec_ValidKeepsValues(t *testing.T) {
	t.Parallel()

	req := trainingdto.RecordTrainingRequest{
		Exercises: []trainingdto.RecordTrainingExerciseRequest{
			{
				ExerciseID:   uuid.NewString(),
				DisplayOrder: 1,
				Sets: []trainingdto.RecordTrainingSetRequest{
					{SetNumber: 1, WeightKg: "9", Reps: 12},
				},
			},
		},
	}

	spec, err := req.ToSpec()
	require.NoError(t, err)
	require.Len(t, spec.Exercises, 1)
	require.Len(t, spec.Exercises[0].Sets, 1)

	set := spec.Exercises[0].Sets[0]
	assert.Equal(t, "9", set.WeightKg.String())
	assert.Equal(t, 12, set.Reps.Value())
}
