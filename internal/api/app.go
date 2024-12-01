package api

import (
	_ "example-go-project/docs"

	"github.com/gin-gonic/gin"

	"example-go-project/internal/handlers"
	"example-go-project/internal/middleware"
	"example-go-project/pkg/config"
)

type Application struct {
	Router      *gin.Engine
	UserHandler *handlers.UserHandler
	AuthHandler *handlers.AuthHandler
	Config      *config.Config
}

// @Description Health check response
type HealthResponse struct {
    Status string `json:"status" example:"ok"`
}

func (app *Application) SetupRoutes() {
	// API version group
	v1 := app.Router.Group("/api/v1")

	// Public routes
	public := v1.Group("")
	{
		// @Summary Health check endpoint
		// @Description Get the API's health status
		// @Tags health
		// @Accept json
		// @Produce json
		// @Success 200 {object} api.HealthResponse
		// @Router /api/v1/health [get]
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, HealthResponse{Status: "ok"})
		})

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
}