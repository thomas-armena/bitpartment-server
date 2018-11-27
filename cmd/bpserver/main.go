package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"github.com/thomas-armena/bitpartment-server/pkg/clockcycle"
	"time"
)


func createTenant(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	fmt.Println(params)
	fmt.Println("name: ",params["name"])
	t := models.Tenant{ ID: 1, Name: params["name"] }
	fmt.Println(t)
}

func main(){
	//Initialize server
	corsObj:=handlers.AllowedOrigins([]string{"*"})

	router := mux.NewRouter()
	router.HandleFunc("/create-tenant/{name}", createTenant).Methods("GET")

	go http.ListenAndServe(":8000", handlers.CORS(corsObj)(router))

	fmt.Println("Listening on port 8000")

	//Initialize clock cycle
	update := make(chan int)
	location, _ := time.LoadLocation("EST")
	startTime := time.Date(time.Now().Year(),time.Now().Month(), time.Now().Day(), 10, 0, 0, 0, location)
	clock := clockcycle.ClockCycle{ StartTime: startTime, Interval: time.Duration(10*time.Second), Frequency: 12, Update: update }
	go clock.Start()
	for {
		fmt.Println(<-clock.Update)
	}
}


