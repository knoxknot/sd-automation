package models

import (
	"context"
	"flag"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Person is exported
type Person struct {
	UUID                    primitive.ObjectID `bson:"_id" json:"uuid,omitempty"`
	Survived                bool               `json:"survived"`
	PassengerClass          int                `json:"passengerClass"`
	Name                    string             `json:"name"`
	Sex                     string             `json:"sex"`
	Age                     int                `json:"age"`
	SiblingsOrSpousesAboard int                `json:"siblingsOrSpousesAboard"`
	ParentsOrChildrenAboard int                `json:"parentsOrChildrenAboard"`
	Fare                    float64            `json:"fare"`
}

// const (
// 	conn = "mongodb://127.0.0.1:27017"
// 	db   = "boarding"
// 	coll = "people"
// )

var (
	conn = flag.String("conn", GetEnvVar("MONGO_CONN_URL", "mongodb://127.0.0.1:27017"), "Mongo Connection String")
	db = flag.String("db", GetEnvVar("DB", "boarding"), "Mongo Database")
	coll = flag.String("coll", GetEnvVar("COLL", "people"), "Mongo Collection")
)

// GetEnvVar is exported
func GetEnvVar(desiredValue, defaultValue string) (value string) {
	value = os.Getenv(desiredValue)
	if value == "" {
		value = defaultValue
	}
	return
}

// GetClient is exported
func GetClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetEnvVar(*conn,"mongodb://127.0.0.1:27017")))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// ConnectDB is exported
func ConnectDB() *mongo.Collection {
	client := GetClient()
	err := client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Failed to Connect to MongoDB!", err)
	} else {
		log.Println("Request operation was successful")
	}

	collection := client.Database(GetEnvVar(*db,"boarding")).Collection(GetEnvVar(*coll,"people"))
	return collection
}