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


func main() {


	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/people", getPeople).Methods("GET")
	api.HandleFunc("/people/{uuid}", getPerson).Methods("GET")
	api.HandleFunc("/people/", addPerson).Methods("POST")
	api.HandleFunc("/people/{uuid}", updatePerson).Methods("PUT")
	api.HandleFunc("/people/{uuid}", deletePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}