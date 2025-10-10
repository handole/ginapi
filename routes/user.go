package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"ginapi/controllers"
	"ginapi/middlewares"
)

func UserRoutes(router *gin.Engine, client *mongo.Client) {
	userCollection := client.Database("ginapi").Collection("users")
	addressCollection := client.Database("ginapi").Collection("addresses")

	userController := controllers.UserController{
		UserCollection:    userCollection,
		AddressCollection: addressCollection,
	}

	userRoutes := router.Group("/users")
	userRoutes.Use(middlewares.AuthMiddleware()) // Apply authentication middleware
	{
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.GET("/", userController.GetUsers)
		userRoutes.GET("/profile", userController.GetProfile)
		userRoutes.GET("/:id/addresses", userController.GetUserAddresses)
		userRoutes.POST("/:id/addresses", userController.AddUserAddresses)
	}
}
