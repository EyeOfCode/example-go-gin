package api

import (
	_ "example-go-project/docs"
	"time"

	"github.com/gin-gonic/gin"

	"example-go-project/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *Application) SetupRoutes() {
	// API version group
	v1 := app.Router.Group("/api/v1")

	// 100 req/1s all route
	v1.Use(middleware.RateLimit(100, time.Minute))

	// Public routes
	public := v1.Group("")
	{
		public.GET("/health", app.helperHandler.HealthCheck)

		// Auth routes
		auth := public.Group("/auth")
		// split rate limit auth only
		auth.Use(middleware.RateLimit(20, time.Minute))
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