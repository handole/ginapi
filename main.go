package main

import (
	_ "ginapi/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ginapi/config"
	"ginapi/routes"
)

func main() {
	// connect to the database
	client := config.ConnectDB()
	// initialize routes gin
	router := gin.Default()
	// swagger
	router.GET("/apidocs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// routes
	routes.UserRoutes(router, client)

	// run the server
	router.Run(":8080")
}

// @title           Go Gin Mongo API
// @version         1.0
// @description     RESTful API with Gin and MongoDB
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /
