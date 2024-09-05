package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionString = "mongodb://localhost:27017"

func DbInstance() *mongo.Client {
	Client_Options := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), Client_Options)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongo Connected Successfully")
	return client

}

var Client *mongo.Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Todo").Collection(collectionName)
	return collection
}
