package api

import (
	helperHandler "example-go-project/internal/handlers/helper"
	userHandler "example-go-project/internal/handlers/user"
	"example-go-project/pkg/config"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Router      *gin.Engine
	helperHandler *helperHandler.HealthHandler
	UserHandler *userHandler.UserHandler
	AuthHandler *userHandler.AuthHandler
	Config      *config.Config
}