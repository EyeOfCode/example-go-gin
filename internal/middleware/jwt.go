package middleware

import (
	"example-go-project/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWT(auth *utils.AuthHandler, role ...utils.Role) gin.HandlerFunc {
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
		claims, err := auth.ValidateToken(token)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		if len(role) > 0 {
			roleSlice := make([]utils.Role, len(claims.Roles))
			for i, r := range claims.Roles {
					roleSlice[i] = utils.Role(r)
			}

			if !utils.IsValidRole(roleSlice, role) {
					utils.SendError(c, http.StatusForbidden, "Insufficient permissions")
					c.Abort()
					return
			}
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}