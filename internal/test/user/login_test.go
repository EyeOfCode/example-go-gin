package test

import (
	"encoding/json"
	"example-go-project/internal/dto"
	"example-go-project/internal/handlers"
	"example-go-project/internal/model"
	"example-go-project/internal/service"
	"example-go-project/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CreateTestContext creates a new Gin context for testing
func CreateTestContext(method, path, jsonBody string) (*gin.Context, *httptest.ResponseRecorder) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create a test context
	c, _ := gin.CreateTestContext(w)

	// Create the request
	req := httptest.NewRequest(method, path, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Set the request to the context
	c.Request = req

	return c, w
}

func TestLogin(t *testing.T) {
	// Setup dependencies
	mockRepo := NewMockUserRepository()
	mockRedis := redis.NewClient(&redis.Options{})
	cfg := &config.Config{
		JWTSecretKey:  "test-secret",
		JWTExpiresIn:  "1h",
		JWTRefreshKey: "test-refresh",
		JWTRefreshIn:  "24h",
	}
	userService := service.NewUserService(mockRepo, mockRedis, cfg)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

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
				user := &model.User{
					ID:        primitive.NewObjectID(),
					Email:     "test@example.com",
					Password:  string(hashedPassword), // "password123" hashed
					Roles:     []string{"user"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"success": true,
				"data": gin.H{
					"access_token":  mock.Anything,
					"refresh_token": mock.Anything,
				},
				"message": "Login successful",
			},
		},
		{
			name: "User Not Found",
			input: dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockUserRepository) {
				m.On("FindOne", mock.Anything, bson.M{"email": "nonexistent@example.com"}).Return(nil, mongo.ErrNoDocuments)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: gin.H{
				"success": false,
				"error":   "Failed to find user",
			},
		},
		{
			name: "Invalid Email Format",
			input: dto.LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": []interface{}{
					"Invalid email format",
				},
			},
		},
		{
			name: "Empty Password",
			input: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			setupMock:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": []interface{}{
					"Password is required",
				},
			},
		},
		{
			name: "Wrong Password",
			input: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(m *MockUserRepository) {
				user := &model.User{
					ID:        primitive.NewObjectID(),
					Email:     "test@example.com",
					Password:  string(hashedPassword), // "password123" hashed
					Roles:     []string{"user"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("FindOne", mock.Anything, bson.M{"email": "test@example.com"}).Return(user, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: gin.H{
				"success": false,
				"error":   "Invalid password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock for this test case
			tt.setupMock(mockRepo)

			// Create handler with the real service (which uses our mock repository)
			handler := handlers.NewUserHandler(userService)

			// Create test context
			jsonData, _ := json.Marshal(tt.input)
			c, w := CreateTestContext("POST", "/auth/login", string(jsonData))

			// Execute the handler
			handler.Login(c)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and verify response
			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// For successful login, verify token structure but not exact values
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, response["success"].(bool))
				assert.Equal(t, "Login successful", response["message"])
				data := response["data"].(map[string]interface{})
				assert.NotEmpty(t, data["access_token"])
				assert.NotEmpty(t, data["refresh_token"])
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}

			// Verify all expectations on mock were met
			mockRepo.AssertExpectations(t)
		})
	}
}
