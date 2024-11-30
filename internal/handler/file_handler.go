package handler

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
)

type FileHandler struct {
    db *gorm.DB
}

func (h *FileHandler) Create(c *gin.Context) {
    userId := c.MustGet("user_id").(uint)
    
    // รับไฟล์
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }

    // สร้างชื่อไฟล์ที่ไม่ซ้ำกัน
    filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
    
    // บันทึกไฟล์
    if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
        c.JSON(500, gin.H{"error": "Could not save file"})
        return
    }

    product := model.Product{
        Name:     c.PostForm("name"),
        Price:    c.PostForm("price"),
        ImageURL: "/uploads/" + filename,
        UserID:   userId,
    }

    if err := h.db.Create(&product).Error; err != nil {
        c.JSON(500, gin.H{"error": "Could not create product"})
        return
    }

    c.JSON(201, product)
}

func (h *FileHandler) List(c *gin.Context) {
    var products []model.Product
    if err := h.db.Find(&products).Error; err != nil {
        c.JSON(500, gin.H{"error": "Could not fetch products"})
        return
    }
    c.JSON(200, products)
}

func (h *FileHandler) Get(c *gin.Context) {
    id := c.Param("id")
    var product model.Product
    
    if err := h.db.First(&product, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Product not found"})
        return
    }
    
    c.JSON(200, product)
}

func (h *FileHandler) Update(c *gin.Context) {
    id := c.Param("id")
    userId := c.MustGet("user_id").(uint)
    
    var product model.Product
    if err := h.db.First(&product, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Product not found"})
        return
    }

    // ตรวจสอบเจ้าของ
    if product.UserID != userId {
        c.JSON(403, gin.H{"error": "Unauthorized"})
        return
    }

    // อัพเดทข้อมูล
    if file, err := c.FormFile("image"); err == nil {
        filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
        if err := c.SaveUploadedFile(file, "uploads/"+filename); err != nil {
            c.JSON(500, gin.H{"error": "Could not save file"})
            return
        }
        product.ImageURL = "/uploads/" + filename
    }

    product.Name = c.PostForm("name")
    product.Price = c.PostForm("price")

    if err := h.db.Save(&product).Error; err != nil {
        c.JSON(500, gin.H{"error": "Could not update product"})
        return
    }

    c.JSON(200, product)
}

func (h *FileHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    userId := c.MustGet("user_id").(uint)

    var product model.Product
    if err := h.db.First(&product, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "Product not found"})
        return
    }

    if product.UserID != userId {
        c.JSON(403, gin.H{"error": "Unauthorized"})
        return
    }

    if err := h.db.Delete(&product).Error; err != nil {
        c.JSON(500, gin.H{"error": "Could not delete product"})
        return
    }

    c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
