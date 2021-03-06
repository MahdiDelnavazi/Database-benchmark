package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdidl/Database-benchmark/Mongo"
	"net/http"
)

func main() {

	// connect to mongodb
	Mongo.MongoConfig()
	//defer Mongo.MongoClient.Disconnect()
	// create one bench document to increment counter in "mongo-increment" endpoint
	Mongo.CreateMongoBench()

	router := gin.Default()
	router.PUT("/mongo-increment", Mongo.BenchIncrement)
	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"ping": "pong"})
	})
	router.Run()

}
