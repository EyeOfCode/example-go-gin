package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"example-go-project/internal/dto"
	"example-go-project/internal/model"
	"example-go-project/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"example-go-project/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository implements UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func init() {
	utils.SetupValidator()
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, query bson.D) (int64, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(int64), args.Error(1)
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		input          dto.LoginRequest
		setupMock      func(*MockUserRepository)
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "Success",
			input: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					ID:       primitive.NewObjectID(),
					Email:    "test@example.com",
					Password: string(hashedPassword),
					Roles:    []string{"user"},
				}
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Credentials",
			input: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(m *MockUserRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: gin.H{
				"success": false,
				"error":   "Invalid email or password",
			},
		},
		{
			name: "User Not Found",
			input: dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "nonexistent@example.com").
					Return(nil, mongo.ErrNoDocuments)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: gin.H{
				"success": false,
				"error":   "Failed to find user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			handler := handlers.NewUserHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request
			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute handler
			handler.Login(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name           string
		input          dto.RegisterRequest
		setupMock      func(*MockUserRepository)
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "Success",
			input: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Password:        "Test123!@#",
				ConfirmPassword: "Test123!@#",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").
					Return(nil, mongo.ErrNoDocuments)
				m.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.Email == "test@example.com" && u.Name == "Test User"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Email Already Exists",
			input: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "existing@example.com",
				Password:        "Test123!@#",
				ConfirmPassword: "Test123!@#",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "existing@example.com").
					Return(&model.User{}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"success": false,
				"error":   "Email already exists",
			},
		},
		{
			name: "Password Mismatch",
			input: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Password:        "Test123!@#",
				ConfirmPassword: "DifferentPass123!@#",
			},
			setupMock:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": []interface{}{"ConfirmPassword must be equal to Password"},
			},
		},
		{
			name: "Invalid Password Format",
			input: dto.RegisterRequest{
				Name:            "Test User",
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": []interface{}{"Password must contain at least one uppercase letter, one number, and one special character"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			handler := handlers.NewUserHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.input)
			c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockUserRepository)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			userID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockUserRepository) {
				user := &model.User{
					ID:    primitive.NewObjectID(),
					Name:  "Test User",
					Email: "test@example.com",
					Roles: []string{"user"},
				}
				m.On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "User Not Found",
			userID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockUserRepository) {
				m.On("FindByID", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, mongo.ErrNoDocuments)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			handler := handlers.NewUserHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Set("userID", tt.userID)
			c.Request = httptest.NewRequest("GET", "/user/profile", nil)

			handler.GetProfile(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
