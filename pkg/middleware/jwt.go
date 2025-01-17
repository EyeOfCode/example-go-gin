package middleware

import (
	"example-go-project/internal/model"
	"example-go-project/internal/service"
	"example-go-project/pkg/config"
	"example-go-project/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	userService *service.UserService
	config      *config.Config
}

func NewAuthMiddleware(userService *service.UserService, config *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		config:      config,
	}
}

// Protected validates JWT token and adds user to context
func (m *AuthMiddleware) Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			utils.SendError(c, http.StatusUnauthorized, "Invalid token format")
			c.Abort()
			return
		}

		token := bearerToken[1]
		auth := utils.NewAuthHandler(m.config.JWTSecretKey, m.config.JWTRefreshKey, m.config.JWTExpiresIn, m.config.JWTRefreshIn)

		claims, err := auth.ValidateToken(token)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		if err := m.userService.ValidateTokenWithRedis(c, token); err != nil {
			utils.SendError(c, http.StatusUnauthorized, "Token is invalid or has been revoked")
			c.Abort()
			return
		}

		user, err := m.userService.FindByID(c, claims.UserID)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("token", token)
		c.Set("claims", claims)
		c.Next()
	}
}

// RequireRoles checks if user has required roles
func (m *AuthMiddleware) RequireRoles(roles ...utils.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			utils.SendError(c, http.StatusUnauthorized, "User not found in context")
			c.Abort()
			return
		}

		userObj, ok := user.(*model.User)
		if !ok {
			utils.SendError(c, http.StatusUnauthorized, "Invalid user type in context")
			c.Abort()
			return
		}

		userRoles := make([]utils.Role, len(userObj.Roles))
		for i, r := range userObj.Roles {
			userRoles[i] = utils.Role(r)
		}

		if !utils.IsValidRole(userRoles, roles) {
			utils.SendError(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromContext retrieves user from context
func GetUserFromContext(c *gin.Context) (*model.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userObj, ok := user.(*model.User)
	return userObj, ok
}
