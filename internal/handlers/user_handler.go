package handlers

import (
	"context"
	"example-go-project/internal/service"
	"example-go-project/pkg/middleware"
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
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary Login endpoint
// @Description Post the API's login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login"
// @Router /auth/login [post]
func (u *UserHandler) Login(c *gin.Context) {
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

	user, err := u.userService.FindByEmail(ctx, req.Email)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Failed to find user")
		return
	}

	if user == nil {
		utils.SendError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	tokenPair, err := u.userService.Login(ctx, req.Password, user)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Invalid password")
		return
	}
	utils.SendSuccess(c, http.StatusOK, tokenPair, "Login successful")
}

// @Summary Register endpoint
// @Description Post the API's register
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration details"
// @Router /auth/register [post]
func (u *UserHandler) Register(c *gin.Context) {
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
	existingUser, err := u.userService.FindByEmail(ctx, req.Email)
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

	if len(req.Roles) == 0 {
		req.Roles = []string{string(utils.UserRole)}
	}

	// Create new user
	req.Password = string(hashedPassword)

	// Save to database
	user, err := u.userService.Create(ctx, &req)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, http.StatusOK, user)
}

// @Summary Profile endpoint
// @Description Get the API's get profile
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Router /user/profile [get]
func (u *UserHandler) GetProfile(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User not found")
		return
	}

	res := gin.H{
		"id":    user.ID.Hex(),
		"name":  user.Name,
		"email": user.Email,
		"roles": user.Roles,
	}

	utils.SendSuccess(c, http.StatusOK, res)
}

// @Summary Update endpoint
// @Description Get the API's update user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param request body dto.UpdateProfileRequest true "User update details"
// @Router /user/profile/{id} [put]
func (u *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	id := c.Param("id")

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

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	user, err := u.userService.FindByID(ctx, objID.Hex())
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	updatedUser, err := u.userService.Update(ctx, &req, user.ID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, http.StatusOK, updatedUser, "Profile updated successfully")
}

// @Summary Delete endpoint
// @Description Get the API's delete user
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Router /user/{id} [delete]
func (u *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User not found")
		return
	}

	if user.ID.Hex() == id {
		utils.SendError(c, http.StatusUnauthorized, "You cannot delete yourself")
		return
	}

	if err := u.userService.Delete(ctx, user.ID); err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, http.StatusOK, nil, "User deleted successfully")
}

// @Summary User list endpoint
// @Description Get the API's user list
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number (default: 1)" default(1)
// @Param pageSize query int false "Page size (default: 10)" default(10)
// @Param name query string false "Filter by user name"
// @Router /user/list [get]
func (u *UserHandler) UserList(c *gin.Context) {
	page, pageSize := utils.PaginationParams(c)

	var filter dto.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid filter parameters")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, total, err := u.userService.FindAll(ctx, filter, page, pageSize)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.CreatePagination(page, pageSize, total, users)
	utils.SendSuccess(c, http.StatusOK, response)
}

// @Summary Refresh endpoint
// @Description Post the API's refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Router /auth/refresh [post]
func (u *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors,
			})
		}

		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokenPair, err := u.userService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	if tokenPair == nil {
		utils.SendError(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	utils.SendSuccess(c, http.StatusOK, tokenPair, "Refresh token successful")
}

// @Summary Logout endpoint
// @Description Post the API's logout
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Header 200 {string} X-Refresh-Token "Refresh token for logout"
// @Param X-Refresh-Token header string true "Refresh token"
// @Router /user/logout [get]
func (u *UserHandler) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get access token from context
	token, ok := c.Get("token")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "Token not found")
		return
	}

	// Get refresh token from header
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		utils.SendError(c, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Call logout service with both tokens
	tokenStr, ok := token.(string)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "Token is not a string")
		return
	}

	if err := u.userService.Logout(ctx, tokenStr, refreshToken); err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, http.StatusOK, nil, "Logout successful")
}
