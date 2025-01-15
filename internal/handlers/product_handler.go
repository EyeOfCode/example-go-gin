package handlers

import (
	"context"
	"example-go-project/internal/dto"
	"example-go-project/internal/model"
	productRepository "example-go-project/internal/repository/product"
	userRepository "example-go-project/internal/repository/user"
	"example-go-project/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductHandler struct {
	productRepo productRepository.ProductRepository
	userRepo    userRepository.UserRepository
}

func NewProductHandler(productRepo productRepository.ProductRepository, userRepo userRepository.UserRepository) *ProductHandler {
	return &ProductHandler{
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// @Summary Create product endpoint
// @Description Post the API's create product
// @Tags product
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateProductRequest true "Product details"
// @Router /product [post]
func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	userID, _ := c.Get("userID")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userIDStr, ok := userID.(string)

	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get user ID")
		return
	}

	user, err := p.userRepo.FindByID(ctx, userIDStr)
	if err != nil && err != mongo.ErrNoDocuments {
		utils.SendError(c, http.StatusInternalServerError, "User not found")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := utils.FormatValidationError(err)
		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors,
			})
			return
		}
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	now := time.Now()
	payload := &model.Product{
		Name:      req.Name,
		Price:     req.Price,
		Stock:     req.Stock,
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = p.productRepo.Create(ctx, payload)

	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create product")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Product created successfully")
}

// @Summary Get products endpoint
// @Description Get the API's get products
// @Tags product
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number (default: 1)" default(1)
// @Param pageSize query int false "Page size (default: 10)" default(10)
// @Param name query string false "Filter by product name"
// @Param price query float64 false "Filter by product price"
// Param stock query int false "Filter by product stock"
// @Param user_id query string false "Filter by product user ID"
// @Router /product [get]
func (p *ProductHandler) GetProducts(c *gin.Context) {
	page, pageSize := utils.PaginationParams(c)

	var filter dto.ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid filter parameters")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoFilter := bson.D{}
	if filter.Name != "" {
		mongoFilter = append(mongoFilter, bson.E{
			Key: "name",
			Value: bson.D{{
				Key:   "$regex",
				Value: primitive.Regex{Pattern: filter.Name, Options: "i"},
			}},
		})
	}

	if filter.Price != nil {
		mongoFilter = append(mongoFilter, bson.E{
			Key: "price",
			Value: bson.D{{
				Key:   "$gte",
				Value: filter.Price,
			}},
		})
	}

	if filter.Stock != nil {
		mongoFilter = append(mongoFilter, bson.E{
			Key: "stock",
			Value: bson.D{{
				Key:   "$gte",
				Value: filter.Stock,
			}},
		})
	}

	if filter.UserId != "" {
		mongoFilter = append(mongoFilter, bson.E{
			Key:   "user_id",
			Value: filter.UserId,
		})
	}

	total, err := p.productRepo.Count(ctx, mongoFilter)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to count users: "+err.Error())
		return
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	products, err := p.productRepo.FindAll(ctx, mongoFilter, opts)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.CreatePagination(page, pageSize, total, products)
	utils.SendSuccess(c, http.StatusOK, response)
}
