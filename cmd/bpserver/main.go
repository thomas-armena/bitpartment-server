package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/models"
)

func runTest(w http.ResponseWriter, r *http.Request){
	fmt.Println("Test called")
}

func createTenant(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	fmt.Println(params)
	fmt.Println("name: ",params["name"])
	t := models.Tenant{ ID: 1, Name: params["name"] }
	fmt.Println(t)
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/test", runTest).Methods("GET")
	router.HandleFunc("/create-tenant/{name}", createTenant).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}


