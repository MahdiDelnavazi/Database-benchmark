package Mongo

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mahdidl/Database-benchmark/Mongo/Entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"net/http"
)

var benchId primitive.ObjectID

func CreateMongoBench() {
	benchBson := Entity.MongoBenchEntity{Name: "mimdl", Counter: 0}
	mongodb := MongoClient.Database(MongoName)
	BenchCollection := mongodb.Collection("Bench")
	result, queryError := BenchCollection.InsertOne(context.Background(), benchBson)
	if queryError != nil {
		fmt.Println(queryError)
		return
	}
	benchId = result.InsertedID.(primitive.ObjectID)
}

func BenchIncrement(ginContext *gin.Context) {

	// this query is fine, but we want to get bad data integrity and fix that with ACID
	//_, queryError := BenchCollection.UpdateOne(MongoContext, bson.M{
	//	"_id": benchId,
	//}, bson.D{
	//	{"$inc", bson.D{{"Counter", 1}}},
	//}, options.Update().SetUpsert(true))
	//if queryError != nil {
	//	log.Println(queryError)
	//	return
	//}

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()

	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := MongoClient.StartSession()
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, "err4")
		panic(err)
		return
	}
	defer session.EndSession(context.TODO())

	err = mongo.WithSession(context.TODO(), session, func(sessionContext mongo.SessionContext) error {
		err = session.StartTransaction(txnOpts)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, "err4")
			return err
		}

		mongodb := MongoClient.Database(MongoName)
		BenchCollection := mongodb.Collection("Bench")

		var BenchFromDb Entity.MongoBenchEntity
		readResult := BenchCollection.FindOne(sessionContext, bson.M{"_id": benchId}).Decode(&BenchFromDb)
		if readResult != nil {
			return readResult
		}

		counter := BenchFromDb.Counter + 1
		result, queryError := BenchCollection.UpdateByID(
			sessionContext,
			BenchFromDb.Id,
			bson.D{
				{"$set", bson.D{{"Counter", counter}}},
			},
		)
		if queryError != nil {
			ginContext.JSON(http.StatusInternalServerError, "err4")
			log.Fatal(queryError)
			return queryError
		}

		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, "err5")
			return err
		}
		if err = session.CommitTransaction(sessionContext); err != nil {
			ginContext.JSON(http.StatusInternalServerError, "err6")
			return err
		}
		log.Println("end transaction")
		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			ginContext.JSON(http.StatusInternalServerError, "err8")
			panic(abortErr)
		}
		ginContext.JSON(http.StatusInternalServerError, "err9")
		panic(err)
	}
}
