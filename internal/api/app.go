package api

import (
	_ "example-go-project/docs"
	"example-go-project/pkg/config"
	"time"

	"github.com/gin-gonic/gin"

	helperHandler "example-go-project/internal/handlers/helper"
	pingHandler "example-go-project/internal/handlers/ping"
	productHandlers "example-go-project/internal/handlers/product"
	uploadHandler "example-go-project/internal/handlers/upload"
	userHandler "example-go-project/internal/handlers/user"
	"example-go-project/internal/middleware"
	"example-go-project/pkg/utils"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	Router      *gin.Engine
	helperHandler *helperHandler.HealthHandler
	UserHandler *userHandler.UserHandler
	AuthHandler *utils.AuthHandler
	PingHandler *pingHandler.PingHandler
	ProductHandler *productHandlers.ProductHandler
	UploadHandler *uploadHandler.UploadHandler
	Config      *config.Config
}

func (app *Application) SetupRoutes() {
	// API version group
	v1 := app.Router.Group("/api/v1")

	authJwt := utils.NewAuthHandler(app.Config.JWTSecretKey, app.Config.JWTExpiresIn)

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

		ping := public.Group("/ping")
		{
			ping.POST("/", app.PingHandler.Ping)
		}
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.JWT(authJwt))
	{
		// User routes
		user := protected.Group("/user")
		{
			user.GET("/profile", app.UserHandler.GetProfile)
			user.PUT("/profile/:id", app.UserHandler.UpdateProfile)
		}
	}

	adminProtected := v1.Group("")
	adminProtected.Use(middleware.JWT(authJwt, utils.AdminRole))
	{
		adminProtected.POST("/local_upload",  app.UploadHandler.UploadMultipleLocalFiles)
		adminProtected.DELETE("/local_upload/:id", app.UploadHandler.DeleteFile)
		adminProtected.GET("/local_upload", app.UploadHandler.GetFileAll)
		
		admin := adminProtected.Group("/user")
		{
			admin.DELETE("/profile/:id", app.UserHandler.DeleteUser)
			admin.GET("/list", app.UserHandler.UserList)
		}
		product := adminProtected.Group("/product")
		{
			product.POST("/", app.ProductHandler.CreateProduct)
			product.GET("/", app.ProductHandler.GetProducts)
		}
	}

	// Setup Swagger
	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}