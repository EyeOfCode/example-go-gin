package handlers

import (
	"context"
	"example-go-project/internal/model"
	userRepository "example-go-project/internal/repository/user"
	"example-go-project/pkg/utils"
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
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        errors := utils.FormatValidationError(err)
        if len(errors) > 0 {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": errors,
            })
            return
        }
        
        utils.SendError(c, http.StatusInternalServerError, err.Error())
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Check if email already exists
    existingUser, err := h.userRepo.FindByEmail(ctx, req.Email)
    if err != nil && err != mongo.ErrNoDocuments {
        utils.SendError(c, http.StatusInternalServerError, "Failed to check existing user")
        return
    }
    if existingUser != nil {
        utils.SendError(c, http.StatusBadRequest, "Email already exists")
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        utils.SendError(c, http.StatusInternalServerError, "Failed to hash password")
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
        utils.SendError(c, http.StatusInternalServerError, err.Error())
        return
    }

    res := dto.RegisterResponse{
        ID:    user.ID.Hex(),
        Name:  user.Name,
        Email: user.Email,
    }
    
    utils.SendSuccess(c, http.StatusOK, res)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... get profile logic ...
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... update profile logic ...
}