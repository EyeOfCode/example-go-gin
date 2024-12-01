package api

import (
	_ "example-go-project/docs"

	"github.com/gin-gonic/gin"

	"example-go-project/internal/handlers"
	"example-go-project/internal/middleware"
	"example-go-project/internal/service"
	"example-go-project/pkg/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	Router      *gin.Engine
	UserHandler *handlers.UserHandler
	AuthHandler *handlers.AuthHandler
	Config      *config.Config
}

func (app *Application) SetupRoutes() {
	// API version group
	v1 := app.Router.Group("/api/v1")

	// Public routes
	public := v1.Group("")
	{
		public.GET("/health", service.HealthCheck)

		// Auth routes
		auth := public.Group("/auth")
		{
			auth.POST("/register", app.UserHandler.Register)
			auth.POST("/login", app.UserHandler.Login)
		}

		// test api
		public.GET("/example", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "feature endpoint"})
		})
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.JWT(app.AuthHandler))
	{
		// User routes
		user := protected.Group("/user")
		{
			user.GET("/profile", app.UserHandler.GetProfile)
			user.PUT("/profile", app.UserHandler.UpdateProfile)
		}
	}

	// Setup Swagger
	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}