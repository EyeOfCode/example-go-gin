package handler

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
)

type JWTAuthHandler struct {
    db *gorm.DB
}

func (h *JWTAuthHandler) Register(c *gin.Context) {
    var user model.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(500, gin.H{"error": "Could not hash password"})
        return
    }

    user.Password = string(hashedPassword)
    if err := h.db.Create(&user).Error; err != nil {
        c.JSON(500, gin.H{"error": "Could not create user"})
        return
    }

    c.JSON(201, gin.H{"message": "User created successfully"})
}

func (h *JWTAuthHandler) Login(c *gin.Context) {
    var loginData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var user model.User
    if err := h.db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte("your-secret-key"))
    if err != nil {
        c.JSON(500, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(200, gin.H{"token": tokenString})
}