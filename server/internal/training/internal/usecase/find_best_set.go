package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// 種目の最高記録(最重量セット)を取得する読み取り UseCase。
// ExerciseRecordQueryService のみを注入する(書き込み Repository は使わない)。
// 実装パターンは usecase/find_routine_by_id.go を参照。
//
//   - FindBestSetInput { UserID, ExerciseID }
//   - FindBestSetOutput { BestSet *trainingdomain.BestSetView }
//   - FindBestSet struct{ exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService }
//   - Execute: QueryService 呼び出し → err は myerror でラップ
//       記録なし(nil)は 404 ではなく「記録なし」として正常応答にする想定
//       (handler 側で 204 / null を返す)
//   - NewFindBestSet(...)

type FindBestSetInput struct {
	UserID     valueobject.UserID
	ExerciseID valueobject.ExerciseID
}

type FindBestSetOutput struct {
	BestSet *trainingdomain.BestSetView
}

type FindBestSet struct {
	exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService
}

func (uc *FindBestSet) Execute(ctx context.Context, input FindBestSetInput) (*FindBestSetOutput, error) {
	bestSet, err := uc.exerciseRecordQueryService.FindBestSet(ctx, input.UserID, input.ExerciseID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	// 記録なし(nil)は 404 ではなく「記録なし」として正常応答にする想定
	if bestSet == nil {
		return nil, nil
	}
	return &FindBestSetOutput{BestSet: bestSet}, nil
}

func NewFindBestSet(exerciseRecordQueryService trainingdomain.ExerciseRecordQueryService) *FindBestSet {
	return &FindBestSet{exerciseRecordQueryService: exerciseRecordQueryService}
}
