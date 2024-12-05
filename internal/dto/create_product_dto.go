package dto

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required,min=3,max=30"`
	Price float64 `json:"price" binding:"required"`
	Stock int     `json:"stock" binding:"required"`
}