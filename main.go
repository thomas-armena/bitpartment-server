package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

func RunTest(w http.ResponseWriter, r *http.Request){
	fmt.Println("Test called")
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/test", RunTest).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}


