package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.com/knoxknot/csproject/application/models"
)

// implement a retrival of the people on the list
func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var people []models.Person
	collection := models.ConnectDB()
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var person models.Person
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

	var person models.Person
	collection := models.ConnectDB()
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
	var person models.Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.UUID = primitive.NewObjectID()

	collection := models.ConnectDB()
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

	collection := models.ConnectDB()

	var oldPerson models.Person
	err = collection.FindOne(context.TODO(), bson.M{"_id": uuid}).Decode(&oldPerson)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	var person models.Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	Survived := r.FormValue("survived")
	PassengerClass := r.FormValue("passengerClass ")
	Name := r.FormValue("name")
	Sex := r.FormValue("sex")
	Age := r.FormValue("age")
	SiblingsOrSpousesAboard := r.FormValue("siblingsOrSpousesAboard")
	ParentsOrChildrenAboard := r.FormValue("parentsOrChildrenAboard")
	Fare := r.FormValue("fare")

	if Survived  == " "  {
		person.Survived  = oldPerson.Survived
	}

	if PassengerClass == " "  {
		person.PassengerClass = oldPerson.PassengerClass
	}

	if Name == ""  {
		person.Name = oldPerson.Name
	}

	if Sex  == " "  {
		person.Sex  = oldPerson.Sex
	}

	if Age == " "  {
		person.Age = oldPerson.Age
	}

	if SiblingsOrSpousesAboard == " "  {
		person.SiblingsOrSpousesAboard = oldPerson.SiblingsOrSpousesAboard
	}

	if ParentsOrChildrenAboard == " "  {
		person.ParentsOrChildrenAboard = oldPerson.ParentsOrChildrenAboard
	}

	if Fare == " "  {
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

	collection := models.ConnectDB()
	objectToDelete, err := collection.DeleteOne(context.TODO(), bson.M{"_id": uuid})
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(objectToDelete.DeletedCount)

}