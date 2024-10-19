package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"golang-techque/database"
	middleware "golang-techque/middleware"
	routes "golang-techque/routes"
)

var foodCollections *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {

	//seting up the port

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	// using router gin

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoicesRoutes(router)

	router.Run(":" + port)

}
