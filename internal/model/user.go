package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string            `bson:"email" json:"email"`
	Password  string            `bson:"password" json:"-"` // "-" means this field won't be included in JSON
	Name      string            `bson:"name" json:"name"`
	Roles			[]string					`bson:"roles" json:"roles,omitempty"`
	Products  []*Product         `bson:"products,omitempty"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

type UserResponseOnProduct struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}