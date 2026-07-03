package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	MongoDB := os.Getenv("MONGODB")
	if MongoDB == "" {
		MongoDB = "mongodb://127.0.0.1:27017"
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDB))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_Name")
	if dbName == "" {
		dbName = os.Getenv("DB_name")
	}
	if dbName == "" {
		dbName = "restaurant_db"
	}

	return client.Database(dbName).Collection(collectionName)
}
