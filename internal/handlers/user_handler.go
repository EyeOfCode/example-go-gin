package handlers

import (
	"example-go-project/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
    return &UserHandler{
        userRepo: userRepo,
    }
}

func (h *UserHandler) Login(c *gin.Context) {
    // ... login logic ...
    // ใช้ h.authService.GenerateToken() เพื่อสร้าง token
}

func (h *UserHandler) Register(c *gin.Context) {
    // ... register logic ...
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... get profile logic ...
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
    // userID, _ := c.Get("userID")
    // ... update profile logic ...
}