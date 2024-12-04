package handlers

import (
	"context"
	"example-go-project/internal/model"
	userRepository "example-go-project/internal/repository/user"
	"example-go-project/pkg/utils"
	"fmt"
	"net/http"
	"os"
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

// @Summary Login endpoint
// @Description Post the API's login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
    var req dto.LoginRequest

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

    user, err := h.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        utils.SendError(c, http.StatusInternalServerError, "Failed to find user")
        return
    }

    if user == nil {
        utils.SendError(c, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        utils.SendError(c, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    auth := utils.NewAuthHandler(os.Getenv("JWT_SECRET"), os.Getenv("JWT_EXPIRY"))
    token, err := auth.GenerateToken(user.ID.Hex())
    if err != nil {
        utils.SendError(c, http.StatusInternalServerError, "Failed to generate token")
        return
    }
    res := gin.H{
        "token": token,
    }
    utils.SendSuccess(c, http.StatusOK, res, "Login successful")
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

    res := gin.H{
        "id":    user.ID.Hex(),
        "name":  user.Name,
        "email": user.Email,
    }
    
    utils.SendSuccess(c, http.StatusOK, res)
}

// @Summary Profile endpoint
// @Description Get the API's get profile
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    userID, _ := c.Get("userID")
    fmt.Println(userID)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    userIDStr, ok := userID.(string)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid user ID format",
        })
        return
    }

    user, err := h.userRepo.FindByID(ctx, userIDStr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "failed to fetch user profile",
        })
        return
    }

    res := gin.H{
        "id":    user.ID.Hex(),
        "name":  user.Name,
        "email": user.Email,
    }
    
    utils.SendSuccess(c, http.StatusOK, res)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... update profile logic ...
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... delete user logic ...
}