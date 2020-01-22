package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

const (
	conn = "mongodb://localhost:27017"
	db   = "boarding"
	coll = "people"
)

func GetClient() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func ConnectDB() *mongo.Collection {
	client := GetClient()
	err := client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Failed to Connect to MongoDB!", err)
	} else {
		log.Println("Request operation was successful")
	}

	collection := client.Database(db).Collection(coll)
	return collection
}

// implement a retrival of the people on the list
func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var people []Person
	collection := ConnectDB()
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(people)
}

// implement a retrival of a person on the list
func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	uuid, err := primitive.ObjectIDFromHex(params["uuid"])
	if err != nil {
		log.Fatal(err)
	}

	var person Person
	collection := ConnectDB()
	err = collection.FindOne(context.TODO(), bson.M{"_id": uuid}).Decode(&person)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(person)
}

// implement the addition of a person to the list
func addPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.UUID = primitive.NewObjectID()

	collection := ConnectDB()
	newPerson, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(newPerson.InsertedID.(primitive.ObjectID))
}

// implement the updating of a particular person on the list
func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	uuid, err := primitive.ObjectIDFromHex(params["uuid"])
	if err != nil {
		log.Fatal(err)
	}

	collection := ConnectDB()

	var oldPerson Person
	err = collection.FindOne(context.TODO(), bson.M{"_id": uuid}).Decode(&oldPerson)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	Survived := r.FormValue("survived")
	PassengerClass := r.FormValue("passengerClass ")
	Name := r.FormValue("name")
	Sex := r.FormValue("sex")
	Age := r.FormValue("age")
	SiblingsOrSpousesAboard := r.FormValue("siblingsOrSpousesAboard")
	ParentsOrChildrenAboard := r.FormValue("parentsOrChildrenAboard")
	Fare := r.FormValue("fare")

	if Survived  == ""  {
		person.Survived  = oldPerson.Survived
	}

	if PassengerClass == ""  {
		person.PassengerClass = oldPerson.PassengerClass
	}

	if Name == ""  {
		person.Name = oldPerson.Name
	}

	if Sex  == ""  {
		person.Sex  = oldPerson.Sex
	}

	if Age == ""  {
		person.Age = oldPerson.Age
	}

	if SiblingsOrSpousesAboard == ""  {
		person.SiblingsOrSpousesAboard = oldPerson.SiblingsOrSpousesAboard
	}

	if ParentsOrChildrenAboard == ""  {
		person.ParentsOrChildrenAboard = oldPerson.ParentsOrChildrenAboard
	}

	if Fare == ""  {
		person.Fare = oldPerson.Fare
	}

	objectDataToUpdate := bson.M{
		"$set": bson.M{
		"survived": person.Survived,
		"passengerClass": person.PassengerClass,
		"name": person.Name,
		"sex": person.Sex,
		"age": person.Age,
		"siblingsOrSpousesAboard": person.SiblingsOrSpousesAboard,
		"parentsOrChildrenAboard": person.ParentsOrChildrenAboard,
		"fare": person.Fare,		
		},
	}

	objectToUpdate, err := collection.UpdateOne(context.TODO(), bson.M{"_id": uuid}, objectDataToUpdate)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(objectToUpdate.ModifiedCount)
}

// implements the deletion of a particular person on the list
func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	uuid, err := primitive.ObjectIDFromHex(params["uuid"])
	if err != nil {
		log.Fatal(err)
	}

	collection := ConnectDB()
	objectToDelete, err := collection.DeleteOne(context.TODO(), bson.M{"_id": uuid})
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(objectToDelete.DeletedCount)

}


func main() {
	client := GetClient()
	err := client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Failed to Connect to MongoDB!", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/people", getPeople).Methods("GET")
	api.HandleFunc("/people/{uuid}", getPerson).Methods("GET")
	api.HandleFunc("/people/", addPerson).Methods("POST")
	api.HandleFunc("/people/{uuid}", updatePerson).Methods("PUT")
	api.HandleFunc("/people/{uuid}", deletePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}