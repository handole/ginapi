package routes

import (
	"ginapi/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegionRoutes(router *gin.Engine, db *mongo.Database) {
	regionController := &controllers.RegionController{
		Collection: db.Collection("region"),
	}

	regionRoutes := router.Group("/regions")
	{
		regionRoutes.POST("", regionController.CreateRegion)
		regionRoutes.GET("", regionController.GetRegions)
		regionRoutes.GET("/:id", regionController.GetRegionByID)
		regionRoutes.PUT("/:id", regionController.UpdateRegion)
		regionRoutes.DELETE("/:id", regionController.DeleteRegion)
	}
}
