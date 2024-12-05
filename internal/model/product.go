package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID        primitive.ObjectID 									`bson:"_id,omitempty" json:"id"`
	Name      string             									`bson:"name" json:"name"`
	Price     float64            									`bson:"price" json:"price"`
	Stock     int                									`bson:"stock" json:"stock"`
	UserID    primitive.ObjectID 									`bson:"user_id"`
  User      *UserResponseOnProduct              `bson:"user,omitempty"`
	CreatedAt time.Time          									`bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          									`bson:"updated_at" json:"updated_at"`
}