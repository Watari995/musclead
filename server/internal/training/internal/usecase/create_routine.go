package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// maxRoutinesForFreeUser は無料ユーザーが作成できるルーティンの最大数
const maxRoutinesForFreeUser = 3

type CreateRoutineInput struct {
	UserID      valueobject.UserID
	RoutineSpec trainingdomain.RoutineSpec
}

type CreateRoutineOutput struct {
	ID valueobject.RoutineID
}

type CreateRoutine struct {
	routineRepo       trainingdomain.RoutineRepository
	subscriptionQuery purchasepublicfunctions.SubscriptionQuery
}

func (uc *CreateRoutine) Execute(ctx context.Context, input CreateRoutineInput) (*CreateRoutineOutput, error) {
	// サブスクチェック
	isPro, err := uc.subscriptionQuery.IsPro(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if !isPro {
		routineCount, err := uc.routineRepo.CountByUserID(ctx, input.UserID)
		if err != nil {
			return nil, myerror.NewInternalError().Wrap(err)
		}
		if routineCount >= maxRoutinesForFreeUser {
			return nil, myerror.NewRoutineLimitReachedError()
		}
	}
	// 末尾に並べるため次の表示順を採番する
	next, err := uc.routineRepo.NextDisplayOrder(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	displayOrder, err := valueobject.NewNonNegativeInt(next)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	// 作成
	routine := trainingdomain.CreateRoutine(input.RoutineSpec, input.UserID, *displayOrder)
	if err := uc.routineRepo.Save(ctx, routine); err != nil {
		if myerror.IsCode(err, myerror.ErrorCodes.Training.RoutineNameAlreadyExistsError) {
			return nil, err
		}
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &CreateRoutineOutput{ID: routine.ID()}, nil
}

func NewCreateRoutine(routineRepo trainingdomain.RoutineRepository, subscriptionQuery purchasepublicfunctions.SubscriptionQuery) *CreateRoutine {
	return &CreateRoutine{routineRepo: routineRepo, subscriptionQuery: subscriptionQuery}
}
