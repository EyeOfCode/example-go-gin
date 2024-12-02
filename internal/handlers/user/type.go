package handlers

import (
	userRepository "example-go-project/internal/repository/user"

	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	secretKey string
	expiresIn string
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type UserHandler struct {
    userRepo userRepository.UserRepository
}

type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}