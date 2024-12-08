package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStorage struct {
	ID        primitive.ObjectID 									`bson:"_id,omitempty" json:"id"`
	Name      string             									`bson:"name" json:"name"`
	Original  string             									`bson:"original" json:"original"`
	BasePath  string             									`bson:"base_path" json:"base_path"`
	Dir       string             									`bson:"url" json:"url"`
	UserID    primitive.ObjectID 									`bson:"user_id"`
	CreatedAt time.Time          									`bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          									`bson:"updated_at" json:"updated_at"`
}