package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "strings"
)

func JWTAuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        // ตัด Bearer ออกจาก token
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

        // ตรวจสอบ token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // เพิ่มข้อมูล claims ใน context
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", uint(claims["user_id"].(float64)))
        }

        c.Next()
    }
}