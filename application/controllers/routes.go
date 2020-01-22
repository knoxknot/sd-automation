package controllers

import (
	"context"
	"flag"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gorilla/mux"
	"gitlab.com/knoxknot/csproject/application/models"
)

var (
	endpoint = flag.String("endpoint", models.GetEnvVar("ENDPOINT", ":8080"), "Server Endpoint")
)

// ServeAPI is exported
func ServeAPI() {
	client := models.GetClient()
	err := client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Failed to Connect to MongoDB!", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	r := mux.NewRouter().StrictSlash(true)
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/people", getPeople).Methods("GET")
	api.HandleFunc("/people/{uuid}", getPerson).Methods("GET")
	api.HandleFunc("/people", addPerson).Methods("POST")
	api.HandleFunc("/people/{uuid}", updatePerson).Methods("PUT")
	api.HandleFunc("/people/{uuid}", deletePerson).Methods("DELETE")

	log.Printf("Server is listening on port %s", *endpoint)
	log.Fatal(http.ListenAndServe(*endpoint, r))
}