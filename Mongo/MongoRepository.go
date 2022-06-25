package Mongo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mahdidl/Database-benchmark/Mongo/Entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var benchId primitive.ObjectID

func CreateMongoBench() {
	benchBson := Entity.MongoBenchEntity{Name: "mimdl", Counter: 0}
	result, queryError := BenchCollection.InsertOne(MongoContext, benchBson)
	if queryError != nil {
		fmt.Println(queryError)
		return
	}
	benchId = result.InsertedID.(primitive.ObjectID)
}

func BenchIncrement(*gin.Context) {
	_, queryError := BenchCollection.UpdateOne(MongoContext, bson.M{
		"_id": benchId,
	}, bson.D{
		{"$inc", bson.D{{"Counter", 1}}},
	}, options.Update().SetUpsert(true))
	if queryError != nil {
		log.Println(queryError)
		return
	}
}
