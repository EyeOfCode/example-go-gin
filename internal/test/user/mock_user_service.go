package test

import (
	"context"
	"example-go-project/internal/dto"
	"example-go-project/internal/model"
	"example-go-project/internal/repository"
	"example-go-project/pkg/config"
	"example-go-project/pkg/utils"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserService is a mock implementation of the UserService
type MockUserService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
	config      *config.Config
	mock.Mock
}

// NewMockUserService creates a new instance of MockUserService
func NewMockUserService(userRepo repository.UserRepository, redisClient *redis.Client, cfg *config.Config) *MockUserService {
	return &MockUserService{
		userRepo:    userRepo,
		redisClient: redisClient,
		config:      cfg,
	}
}

func (m *MockUserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Create(ctx context.Context, req *dto.RegisterRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, password string, user *model.User) (*utils.TokenPair, error) {
	args := m.Called(ctx, password, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.TokenPair), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, payload *dto.UpdateProfileRequest, id primitive.ObjectID) (*model.User, error) {
	args := m.Called(ctx, payload, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) FindAll(ctx context.Context, filter dto.UserFilter, page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) RefreshToken(ctx context.Context, refreshToken string) (*utils.TokenPair, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*utils.TokenPair), args.Error(1)
}

func (m *MockUserService) ValidateTokenWithRedis(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	args := m.Called(ctx, accessToken, refreshToken)
	return args.Error(0)
}
