package trainingdomain

import (
	"context"
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

// 種目(exercise)の記録/成績系の読み取り(CQRS の Query 側)。
// 集約 Repository とは分離し、読み取り最適化した View を返す。
// 実装パターンは routine_query.go を参照。
//
//   - BestSetView: FindBestSet の read DTO(weightKg / reps / performedAt / trainingID など)
//       ※ View はメソッド固有名にする(将来の読み取りは別の形を返しうるため)
//   - ExerciseRecordQueryService interface:
//       FindBestSet(ctx, userID, exerciseID) (*BestSetView, error)
//         → 同一種目の全セットから weight DESC, reps DESC で先頭1件
//         → 記録なしは (nil, nil)
//   - 将来: FindAllTimeMax / FindLastSessionBest など同系統の読み取りをここに追加
//          (それぞれ専用 View を返す)

type BestSetView struct {
	WeightKg    valueobject.NonNegativeDecimal
	Reps        valueobject.NonNegativeInt
	PerformedAt time.Time
	TrainingID  valueobject.TrainingID
}

type ExerciseRecordQueryService interface {
	FindBestSet(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID) (*BestSetView, error)
}
