package service

import (
	"context"
	"example-go-project/internal/dto"
	"example-go-project/internal/model"
	"example-go-project/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, payload *dto.CreateProductRequest, userId primitive.ObjectID) (*model.Product, error) {
	now := time.Now()
	req := &model.Product{
		Name:     payload.Name,
		Price:    payload.Price,
		Stock:    payload.Stock,
		UserID:   userId,
		CreatedAt: now,
		UpdatedAt: now,
	}
	res, err := p.productRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ProductService) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]*model.Product, error) {
	products, err := p.productRepo.FindAll(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductService) Count(ctx context.Context, query bson.D) (int64, error) {
	return p.productRepo.Count(ctx, query)
}