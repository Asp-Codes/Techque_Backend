package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	MongoDb := "mongodb://localhost:27017" //linking the local database.
	fmt.Println(MongoDb)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //giving 10 secons as a timeout

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDb)) //adding the context to the local database.
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	err = client.Ping(ctx, nil) //calling the connection

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	return client

}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("techque").Collection(collectionName)
	//collection name is used to get a specific collection which we need to work with.(probably tables u can say).
	return collection
}
