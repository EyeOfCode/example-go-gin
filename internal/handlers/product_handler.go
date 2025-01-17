package handlers

import (
	"context"
	"example-go-project/internal/dto"
	"example-go-project/internal/service"
	"example-go-project/pkg/middleware"
	"example-go-project/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductHandler struct {
	productService *service.ProductService
	userService    *service.UserService
}

func NewProductHandler(productService *service.ProductService, userService *service.UserService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		userService:    userService,
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User not found")
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
	res, err := p.productService.CreateProduct(ctx, &req, user.ID)

	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create product")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, res, "Product created successfully")
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

	total, err := p.productService.Count(ctx, mongoFilter)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to count users: "+err.Error())
		return
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	products, err := p.productService.FindAll(ctx, mongoFilter, opts)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := utils.CreatePagination(page, pageSize, total, products)
	utils.SendSuccess(c, http.StatusOK, response)
}
