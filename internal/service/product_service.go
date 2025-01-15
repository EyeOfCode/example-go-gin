package service

import (
	"example-go-project/internal/repository"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// func (p *ProductService) CreateProduct(ctx context.Context, product *repository.ProductRepository) error {
// 	return p.productRepo.Create(ctx, product)
// }

// func (p *ProductService) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]repository.ProductRepository, error) {
// 	return p.productRepo.FindAll(ctx, query, opts)
// }