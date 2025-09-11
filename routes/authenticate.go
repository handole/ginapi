package routes

import (
	"ginapi/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthRoutes(router *gin.Engine, db *mongo.Database) {
	collection := db.Collection("users")
	authController := controllers.AuthController{Collection: collection}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
	}
}
