package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ReorderRoutinesInput struct {
	UserID     valueobject.UserID
	OrderedIDs []valueobject.RoutineID
}

type ReorderRoutines struct {
	routineRepo trainingdomain.RoutineRepository
	txManager   dbtx.TransactionManager
}

func (uc *ReorderRoutines) Execute(ctx context.Context, input ReorderRoutinesInput) error {
	routines, err := uc.routineRepo.FindAllByUserID(ctx, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	// 渡された並び順が、 ユーザーの保有するルーティン全件のちょうど 1 度ずつの並べ替えであることを検証する。
	byID := make(map[string]*trainingdomain.Routine, len(routines))
	for _, r := range routines {
		byID[r.ID().Value()] = r
	}
	if len(input.OrderedIDs) != len(byID) {
		return myerror.NewBadRequestError().SetMessage("ordered ids must contain every routine exactly once")
	}
	seen := make(map[string]struct{}, len(input.OrderedIDs))
	for _, id := range input.OrderedIDs {
		key := id.Value()
		if _, ok := byID[key]; !ok {
			return myerror.NewBadRequestError().SetMessage("ordered ids contain an unknown routine")
		}
		if _, dup := seen[key]; dup {
			return myerror.NewBadRequestError().SetMessage("ordered ids contain a duplicate routine")
		}
		seen[key] = struct{}{}
	}

	return uc.txManager.Processing(ctx, func(txCtx context.Context) error {
		for index, id := range input.OrderedIDs {
			displayOrder, err := valueobject.NewNonNegativeInt(index)
			if err != nil {
				return myerror.NewInternalError().Wrap(err)
			}
			routine := byID[id.Value()]
			routine.SetDisplayOrder(*displayOrder)
			if err := uc.routineRepo.Save(txCtx, routine); err != nil {
				return myerror.NewInternalError().Wrap(err)
			}
		}
		return nil
	})
}

func NewReorderRoutines(routineRepo trainingdomain.RoutineRepository, txManager dbtx.TransactionManager) *ReorderRoutines {
	return &ReorderRoutines{routineRepo: routineRepo, txManager: txManager}
}
