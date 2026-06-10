package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ReorderExercisesInput struct {
	UserID     valueobject.UserID
	OrderedIDs []valueobject.ExerciseID
}

type ReorderExercises struct {
	exerciseRepo trainingdomain.ExerciseRepository
	txManager    dbtx.TransactionManager
}

func (uc *ReorderExercises) Execute(ctx context.Context, input ReorderExercisesInput) error {
	exercises, err := uc.exerciseRepo.FindAllByUserID(ctx, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	// 渡された並び順が、 ユーザーの保有する種目全件のちょうど 1 度ずつの並べ替えであることを検証する。
	byID := make(map[string]*trainingdomain.Exercise, len(exercises))
	for _, e := range exercises {
		byID[e.ID().Value()] = e
	}
	if len(input.OrderedIDs) != len(byID) {
		return myerror.NewBadRequestError().SetMessage("ordered ids must contain every exercise exactly once")
	}
	seen := make(map[string]struct{}, len(input.OrderedIDs))
	for _, id := range input.OrderedIDs {
		key := id.Value()
		if _, ok := byID[key]; !ok {
			return myerror.NewBadRequestError().SetMessage("ordered ids contain an unknown exercise")
		}
		if _, dup := seen[key]; dup {
			return myerror.NewBadRequestError().SetMessage("ordered ids contain a duplicate exercise")
		}
		seen[key] = struct{}{}
	}

	return uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		for index, id := range input.OrderedIDs {
			displayOrder, err := valueobject.NewNonNegativeInt(index)
			if err != nil {
				return myerror.NewInternalError().Wrap(err)
			}
			exercise := byID[id.Value()]
			exercise.SetDisplayOrder(*displayOrder)
			if err := uc.exerciseRepo.Save(txCtx, exercise); err != nil {
				return myerror.NewInternalError().Wrap(err)
			}
		}
		return nil
	})
}

func NewReorderExercises(exerciseRepo trainingdomain.ExerciseRepository, txManager dbtx.TransactionManager) *ReorderExercises {
	return &ReorderExercises{exerciseRepo: exerciseRepo, txManager: txManager}
}
