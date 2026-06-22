package trainingusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Watari995/musclead/internal/myerror"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRoutine_ProUser_Success(t *testing.T) {
	t.Parallel()
	routineRepo := new(MockRoutineRepository)
	subQuery := new(MockSubscriptionQuery)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Push Day")

	subQuery.On("IsPro", mock.Anything, mock.Anything).Return(true, nil)
	routineRepo.On("NextDisplayOrder", mock.Anything, mock.Anything).Return(0, nil)
	routineRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	uc := trainingusecase.NewCreateRoutine(routineRepo, subQuery)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateRoutineInput{
		UserID:      userID,
		RoutineSpec: trainingdomain.RoutineSpec{Name: *name},
	})

	assert.NoError(t, err)
	assert.NotNil(t, out)
	routineRepo.AssertExpectations(t)
	subQuery.AssertExpectations(t)
}

func TestCreateRoutine_FreeUser_UnderLimit_Success(t *testing.T) {
	t.Parallel()
	routineRepo := new(MockRoutineRepository)
	subQuery := new(MockSubscriptionQuery)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Push Day")

	subQuery.On("IsPro", mock.Anything, mock.Anything).Return(false, nil)
	routineRepo.On("CountByUserID", mock.Anything, mock.Anything).Return(2, nil)
	routineRepo.On("NextDisplayOrder", mock.Anything, mock.Anything).Return(2, nil)
	routineRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	uc := trainingusecase.NewCreateRoutine(routineRepo, subQuery)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateRoutineInput{
		UserID:      userID,
		RoutineSpec: trainingdomain.RoutineSpec{Name: *name},
	})

	assert.NoError(t, err)
	assert.NotNil(t, out)
}

func TestCreateRoutine_FreeUser_LimitReached(t *testing.T) {
	t.Parallel()
	routineRepo := new(MockRoutineRepository)
	subQuery := new(MockSubscriptionQuery)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Pull Day")

	subQuery.On("IsPro", mock.Anything, mock.Anything).Return(false, nil)
	routineRepo.On("CountByUserID", mock.Anything, mock.Anything).Return(3, nil)

	uc := trainingusecase.NewCreateRoutine(routineRepo, subQuery)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateRoutineInput{
		UserID:      userID,
		RoutineSpec: trainingdomain.RoutineSpec{Name: *name},
	})

	assert.Error(t, err)
	assert.Nil(t, out)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Training.RoutineLimitReachedError))
}

func TestCreateRoutine_NameAlreadyExists(t *testing.T) {
	t.Parallel()
	routineRepo := new(MockRoutineRepository)
	subQuery := new(MockSubscriptionQuery)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Leg Day")

	subQuery.On("IsPro", mock.Anything, mock.Anything).Return(true, nil)
	routineRepo.On("NextDisplayOrder", mock.Anything, mock.Anything).Return(0, nil)
	routineRepo.On("Save", mock.Anything, mock.Anything).Return(myerror.NewRoutineNameAlreadyExistsError())

	uc := trainingusecase.NewCreateRoutine(routineRepo, subQuery)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateRoutineInput{
		UserID:      userID,
		RoutineSpec: trainingdomain.RoutineSpec{Name: *name},
	})

	assert.Error(t, err)
	assert.Nil(t, out)
	assert.True(t, myerror.IsCode(err, myerror.ErrorCodes.Training.RoutineNameAlreadyExistsError))
}

func TestCreateRoutine_SubscriptionQueryError(t *testing.T) {
	t.Parallel()
	routineRepo := new(MockRoutineRepository)
	subQuery := new(MockSubscriptionQuery)
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	name, _ := valueobject.NewString50("Leg Day")

	subQuery.On("IsPro", mock.Anything, mock.Anything).Return(false, errors.New("rpc error"))

	uc := trainingusecase.NewCreateRoutine(routineRepo, subQuery)
	out, err := uc.Execute(context.Background(), trainingusecase.CreateRoutineInput{
		UserID:      userID,
		RoutineSpec: trainingdomain.RoutineSpec{Name: *name},
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}
