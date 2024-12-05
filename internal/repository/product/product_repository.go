package repository

import (
	"context"
	"example-go-project/internal/model"
	repository "example-go-project/internal/repository/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.Product, error)
	Count(ctx context.Context, query bson.D) (int64, error)
}

type productRepository struct {
	collection *mongo.Collection
	user repository.UserRepository
}

func NewProductRepository(db *mongo.Database, user repository.UserRepository) ProductRepository {
	return &productRepository{
		collection: db.Collection("products"),
		user: user,
	}
}

func (p *productRepository) Create(ctx context.Context, product *model.Product) error {
	_, err := p.collection.InsertOne(ctx, product)
	return err
}

func (p *productRepository) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.Product, error) {
	cursor, err := p.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	for i := range products {
			user, err := p.user.FindByID(ctx, products[i].UserID.Hex())
			if err != nil {
					continue
			}
			products[i].User = &model.UserResponseOnProduct{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			}
	}

	return products, nil
}

func (p *productRepository) Count(ctx context.Context, query bson.D) (int64, error) {
	return p.collection.CountDocuments(ctx, query)
}