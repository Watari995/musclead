package trainingusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Watari995/musclead/internal/pagination"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/mock"
)

// --- TrainingRepository ---

type MockTrainingRepository struct{ mock.Mock }

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

// --- ExerciseRepository ---

type MockExerciseRepository struct{ mock.Mock }

func (m *MockExerciseRepository) FindByIDAndUserID(ctx context.Context, id valueobject.ExerciseID, userID valueobject.UserID) (*trainingdomain.Exercise, error) {
	args := m.Called(ctx, id, userID)
	ex, _ := args.Get(0).(*trainingdomain.Exercise)
	return ex, args.Error(1)
}

func (m *MockExerciseRepository) FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*trainingdomain.Exercise, error) {
	args := m.Called(ctx, userID)
	exs, _ := args.Get(0).([]*trainingdomain.Exercise)
	return exs, args.Error(1)
}

func (m *MockExerciseRepository) FindAllByUserIDWithOffsetPagination(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]*trainingdomain.Exercise, pagination.OffsetPaginator, error) {
	args := m.Called(ctx, userID, limit, offset)
	exs, _ := args.Get(0).([]*trainingdomain.Exercise)
	pg, _ := args.Get(1).(pagination.OffsetPaginator)
	return exs, pg, args.Error(2)
}

func (m *MockExerciseRepository) NextDisplayOrder(ctx context.Context, userID valueobject.UserID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockExerciseRepository) Save(ctx context.Context, ex *trainingdomain.Exercise) error {
	return m.Called(ctx, ex).Error(0)
}

func (m *MockExerciseRepository) DeleteByID(ctx context.Context, id valueobject.ExerciseID) error {
	return m.Called(ctx, id).Error(0)
}

// --- RoutineRepository ---

type MockRoutineRepository struct{ mock.Mock }

func (m *MockRoutineRepository) FindByIDAndUserID(ctx context.Context, id valueobject.RoutineID, userID valueobject.UserID) (*trainingdomain.Routine, error) {
	args := m.Called(ctx, id, userID)
	r, _ := args.Get(0).(*trainingdomain.Routine)
	return r, args.Error(1)
}

func (m *MockRoutineRepository) FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*trainingdomain.Routine, error) {
	args := m.Called(ctx, userID)
	rs, _ := args.Get(0).([]*trainingdomain.Routine)
	return rs, args.Error(1)
}

func (m *MockRoutineRepository) CountByUserID(ctx context.Context, userID valueobject.UserID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockRoutineRepository) NextDisplayOrder(ctx context.Context, userID valueobject.UserID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockRoutineRepository) Save(ctx context.Context, r *trainingdomain.Routine) error {
	return m.Called(ctx, r).Error(0)
}

func (m *MockRoutineRepository) DeleteByID(ctx context.Context, id valueobject.RoutineID) error {
	return m.Called(ctx, id).Error(0)
}

// --- ExerciseBestSetTimeseriesCache ---

type MockExerciseBestSetTimeseriesCache struct{ mock.Mock }

func (m *MockExerciseBestSetTimeseriesCache) FindByPeriod(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID, from, to time.Time) ([]*trainingdomain.BestSetView, bool, error) {
	args := m.Called(ctx, userID, exerciseID, from, to)
	views, _ := args.Get(0).([]*trainingdomain.BestSetView)
	return views, args.Bool(1), args.Error(2)
}

func (m *MockExerciseBestSetTimeseriesCache) Save(ctx context.Context, userID valueobject.UserID, bestSet *trainingdomain.BestSetView) error {
	return m.Called(ctx, userID, bestSet).Error(0)
}

func (m *MockExerciseBestSetTimeseriesCache) Evict(ctx context.Context, userID valueobject.UserID, exerciseID valueobject.ExerciseID) error {
	return m.Called(ctx, userID, exerciseID).Error(0)
}

// --- SubscriptionQuery ---

type MockSubscriptionQuery struct{ mock.Mock }

func (m *MockSubscriptionQuery) IsPro(ctx context.Context, userID valueobject.UserID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

// --- fakeTxManager ---

type fakeTxManager struct{}

func (fakeTxManager) Processing(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

// --- helpers ---

func newDummyTraining(t *testing.T) *trainingdomain.Training {
	t.Helper()
	id := valueobject.NewPrimaryID[valueobject.TrainingID]()
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	started := time.Now()
	ended := started.Add(time.Hour)
	memo, _ := valueobject.NewString1000("dummy memo")
	return trainingdomain.NewTraining(id, userID, started, &ended, memo, time.Now(), time.Now(), nil)
}

func newDummyExercise(t *testing.T, userID valueobject.UserID) *trainingdomain.Exercise {
	t.Helper()
	name, _ := valueobject.NewString50("Bench Press")
	order, _ := valueobject.NewNonNegativeInt(0)
	return trainingdomain.CreateExercise(userID, *name, *order)
}

func newDummyRoutine(t *testing.T, userID valueobject.UserID) *trainingdomain.Routine {
	t.Helper()
	name, _ := valueobject.NewString50("Push Day")
	order, _ := valueobject.NewNonNegativeInt(0)
	spec := trainingdomain.RoutineSpec{Name: *name}
	return trainingdomain.CreateRoutine(spec, userID, *order)
}
