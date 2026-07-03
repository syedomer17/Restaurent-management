package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(500, gin.H{"error": "Error occurred while listing orders items"})
			return
		}

		var allOrders []bson.M

		if err = result.All(ctx, &allOrders); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, allOrders)
		defer cancel()
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{
			"order_id": orderId,
		}).Decode(&order)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occured while fetching the order",
			})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if validationErr := validate.Struct(order); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Check if table exists
		err := tableCollection.FindOne(ctx, bson.M{
			"table_id": order.Table_id,
		}).Decode(&table)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Table not found",
			})
			return
		}

		now := time.Now()

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		order.CreatedAt = now
		order.UpdatedAt = now

		result, err := orderCollection.InsertOne(ctx, order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create order",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Order created successfully",
			"result":  result,
		})
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order

		var updateObj primitive.D

		orderId := c.Param("order_id")

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if order.Table_id != nil {
			filter := bson.M{"table_id": order.Table_id}
			err := menuCollection.FindOne(context.TODO(), filter).Decode(&table)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Table not found",
				})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.Table_id})
		}

		upsert := true

		filter := bson.M{"order_id": orderId}

		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&options,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Order update failed",
			})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func DeleteOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderID := c.Param("order_id")

		result, err := orderCollection.DeleteOne(ctx, bson.M{
			"order_id": orderID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete order",
			})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Order deleted successfully",
		})
	}
}

func OrderItemOrderCreator(order models.Order) string {
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	return order.Order_id
}