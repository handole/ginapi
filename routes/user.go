package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"ginapi/controllers"
)

func UserRoutes(router *gin.Engine, client *mongo.Client) {
	collection := client.Database("ginapi").Collection("users")
	userController := controllers.UserController{Collection: collection}

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.GET("/", userController.GetUsers)
	}
}
