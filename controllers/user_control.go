package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	collection *mongo.Collection
}

// create a new user controller
