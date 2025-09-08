package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	createdAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	updatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
