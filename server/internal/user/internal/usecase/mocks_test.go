package userusecase_test

import (
	"context"

	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository は userdomain.UserRepository の偽実装(テスト用)。
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*userdomain.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(*userdomain.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*userdomain.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*userdomain.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, user *userdomain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUserIDs(ctx context.Context) ([]valueobject.UserID, error) {
	args := m.Called(ctx)
	ids, _ := args.Get(0).([]valueobject.UserID)
	return ids, args.Error(1)
}

// MockPasswordHasher は userdomain.PasswordHasher の偽実装(テスト用)。
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(rawPassword string) (*valueobject.HashedPassword, error) {
	args := m.Called(rawPassword)
	hash, _ := args.Get(0).(*valueobject.HashedPassword)
	return hash, args.Error(1)
}

func (m *MockPasswordHasher) Compare(rawPassword string, hash *valueobject.HashedPassword) error {
	args := m.Called(rawPassword, hash)
	return args.Error(0)
}
