package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"

	"golang-restaurant-management/database"
	middleware "golang-restaurant-management/middleware"
	routes "golang-restaurant-management/routes"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment")
	}

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
	routes.PaymentRoutes(router)
	router.Run(":" + port)
}
