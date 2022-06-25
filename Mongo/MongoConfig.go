package Mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	BenchCollection *mongo.Collection
	MongoContext    = context.TODO()

	MongoUrl  = "mongodb://localhost:27017/"
	MongoName = "database-benchmark"
)

// MongoConfig is a func for configuration mongodb
func MongoConfig() {
	client, context, _, err := connect(MongoUrl)
	if err != nil {
		log.Fatal(err)
	}
	mongodb := client.Database(MongoName)

	BenchCollection = mongodb.Collection("Bench")

	ping(client, context)
}

// connect is a func fot connect to mongodb
func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	context, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(context, options.Client().ApplyURI(uri))
	return client, context, cancel, err
}

// ping is a func for make sure we connect to mongodb successfully
func ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}
