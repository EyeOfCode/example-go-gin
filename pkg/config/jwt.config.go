package config

import "os"

type JwtConfig struct {
    SecretKey string
}

func LoadJwtConfig() *JwtConfig {
    return &JwtConfig{
        SecretKey: os.Getenv("JWT_SECRET_KEY"),
    }
}