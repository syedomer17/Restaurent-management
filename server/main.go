package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"golang-restaurant-management/database"
	"golang-restaurant-management/routes"
	"golang-restaurant-management/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New() 
	router.Use(gin.Logger())
	routes.UserRouter(router)
	router.Use(middleware.Authentication())

	routes.FoodRouter(router)
	routes.MenuRouter(router)
	routes.TableRouter(router)
	routes.OrderRouter(router)
	routes.OrderItemRouter(router)
	routes.InvoiceRouter(router)

	router.Run(":" + port)
}