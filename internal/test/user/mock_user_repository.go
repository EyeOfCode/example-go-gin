package test

import (
	"context"
	"example-go-project/internal/model"
	"example-go-project/internal/repository"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements UserRepository interface
var _ repository.UserRepository = &MockUserRepository{}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, payload bson.M, id primitive.ObjectID) (*model.User, error) {
	args := m.Called(ctx, payload, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindOne(ctx context.Context, query bson.M) (*model.User, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.User, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, query bson.D) (int64, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(int64), args.Error(1)
}
