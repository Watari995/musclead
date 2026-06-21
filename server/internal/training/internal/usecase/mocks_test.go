package trainingusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Wataris995/musclead/internal/pagination"
	trainingdomain "github.com/Wataris995/musclead/internal/training/internal/domain"
	"github.com/Wataris995/musclead/internal/valueobject"
	"github.com/stretchr/testify/mock"
)

type MockTrainingRepository struct {
	mock.Mock
}

func (m *MockTrainingRepository) FindByIDAndUserID(ctx context.Context, id valueobject.TrainingID, userID valueobject.UserID) (*trainingdomain.Training, error) {
	args := m.Called(ctx, id, userID)
	t, _ := args.Get(0).(*trainingdomain.Training)
	return t, args.Error(1)
}

func (m *MockTrainingRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]*trainingdomain.Training, pagination.OffsetPaginator, error) {
	args := m.Called(ctx, userID, limit, offset)
	ts, _ := args.Get(0).([]*trainingdomain.Training)
	pg, _ := args.Get(1).(pagination.OffsetPaginator)
	return ts, pg, args.Error(2)
}

func (m *MockTrainingRepository) Save(ctx context.Context, t *trainingdomain.Training) error {
	return m.Called(ctx, t).Error(0)
}

func (m *MockTrainingRepository) DeleteByID(ctx context.Context, id valueobject.TrainingID) error {
	return m.Called(ctx, id).Error(0)
}

type fakeTxManager struct{}

func (fakeTxManager) Processing(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

func newDummyTraining(t *testing.T) *trainingdomain.Training {
	t.Helper()
	id := valueobject.NewPrimaryID[valueobject.TrainingID]()
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	started := time.Now()
	ended := started.Add(time.Hour)
	memo, _ := valueobject.NewString1000("dummy memo")
	return trainingdomain.NewTraining(id, userID, started, &ended, memo, time.Now(), time.Now(), nil)
}
