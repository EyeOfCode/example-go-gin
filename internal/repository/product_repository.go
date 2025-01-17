package repository

import (
	"context"
	"example-go-project/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]*model.Product, error)
	Count(ctx context.Context, query bson.D) (int64, error)
}

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepository{
		collection: db.Collection("products"),
	}
}

func (p *productRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	res, err := p.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}
	productId := res.InsertedID.(primitive.ObjectID)
	product.ID = productId
	return product, nil
}

func (p *productRepository) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]*model.Product, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: query}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: opts.Skip}},
		{{Key: "$limit", Value: opts.Limit}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: "$user"}},
	}

	cursor, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (p *productRepository) Count(ctx context.Context, query bson.D) (int64, error) {
	return p.collection.CountDocuments(ctx, query)
}
