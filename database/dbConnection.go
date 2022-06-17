package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
)

const DbUrl string = "mongodb://127.0.0.1:27017"

func getClient() *mongo.Client {
	// Create client
	client, err := mongo.NewClient(options.Client().ApplyURI(DbUrl))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var Client *mongo.Client = getClient()

func Connect(client *mongo.Client, collectionName string) *mongo.Collection {

	// Create connect
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection(collectionName)

	return collection
}

//func Disconnect(client *mongo.Client) {
//	err := client.Disconnect(context.TODO())
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Connection to MongoDB closed.")
//}
