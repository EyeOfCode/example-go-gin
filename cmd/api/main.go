package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/gorm"
    "example-go-project/pkg/config"
    "example-go-project/pkg/db"
    "example-go-project/internal/middleware"
)

func main() {
    // โหลด config
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // โหลด config
    dbConfig := config.LoadDbConfig()
    jwtConfig := config.LoadJwtConfig()
    
    // เชื่อมต่อ database
    database, err := database.NewPostgresDB(dbConfig.GetDSN())
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // สร้าง Gin router
    r := gin.Default()
    
    // สร้างโฟลเดอร์สำหรับเก็บไฟล์
    os.MkdirAll("uploads", 0755)
    
    // Serve static files
    r.Static("/uploads", "./uploads")

    // Initialize handlers
    authHandler := &JWTAuthHandler{db: database, jwtConfig: jwtConfig}
    productHandler := &FileHandler{db: database}

    // Public routes
    r.POST("/register", authHandler.Register)
    r.POST("/login", authHandler.Login)

    // Protected routes
    api := r.Group("/api")
    api.Use(middleware.JWTAuthMiddleware())
    {
        api.POST("/products", productHandler.Create)
        api.GET("/products", productHandler.List)
        api.GET("/products/:id", productHandler.Get)
        api.PUT("/products/:id", productHandler.Update)
        api.DELETE("/products/:id", productHandler.Delete)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}