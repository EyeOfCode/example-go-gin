package handlers

import (
	"context"
	"example-go-project/internal/model"
	userRepository "example-go-project/internal/repository/user"
	"net/http"
	"time"

	"example-go-project/internal/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
    userRepo userRepository.UserRepository
}

func NewUserHandler(userRepo userRepository.UserRepository) *UserHandler {
    return &UserHandler{
        userRepo: userRepo,
    }
}

func (h *UserHandler) Login(c *gin.Context) {
    // ... login logic ...
    // ใช้ h.authService.GenerateToken() เพื่อสร้าง token
}

// @Summary Register endpoint
// @Description Post the API's register
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Success 200 {object} dto.RegisterResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check if email already exists
    existingUser, err := h.userRepo.FindByEmail(ctx, req.Email)
    if err != nil && err != mongo.ErrNoDocuments {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to check existing user",
        })
        return
    }
    if existingUser != nil {
        c.JSON(http.StatusConflict, gin.H{
            "error": "Email already registered",
        })
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to process password",
        })
        return
    }

    // Create new user
    now := time.Now()
    user := &model.User{
        ID:        primitive.NewObjectID(),
        Name:      req.Name,
        Email:     req.Email,
        Password:  string(hashedPassword),
        CreatedAt: now,
        UpdatedAt: now,
    }

    // Save to database
    if err := h.userRepo.Create(ctx, user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create user",
        })
        return
    }

    // Return success response
    c.JSON(http.StatusCreated, dto.RegisterResponse{
        ID:    user.ID.Hex(),
        Name:  user.Name,
        Email: user.Email,
    })
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... get profile logic ...
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... update profile logic ...
}