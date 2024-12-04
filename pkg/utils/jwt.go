package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	secretKey string
	expiresIn string
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Roles  []string
	jwt.StandardClaims
}

func NewAuthHandler(secretKey, expiresIn string) *AuthHandler {
	return &AuthHandler{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

func (s *AuthHandler) GenerateToken(userID string, roles []string) (string, error) {
	expDuration, _ := time.ParseDuration(s.expiresIn)
	claims := JWTClaims{
		UserID: userID,
		Roles:  roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthHandler) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}