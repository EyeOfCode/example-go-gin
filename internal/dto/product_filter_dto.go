package dto

type ProductFilter struct {
	Name   string   `form:"name"`
	Price  *float64 `form:"price"`
	Stock  *int     `form:"stock"`
	UserId string   `form:"user_id"`
}