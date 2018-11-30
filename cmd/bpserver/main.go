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


	//Define rooms
	bedroom := models.Room{ Name: "bedroom" }
	work := models.Room{ Name: "shop" }

	//Define doings
	sleeping := models.Doing{ Name: "sleeping" }
	working := models.Doing{ Name: "coding" }

	//Define actions
	goToWork := models.Action{ Interval: 32, Doing: &working, Location: &work }
	goToSleep := models.Action{ Interval: 88, Doing: &sleeping, Location: &bedroom }

	//Define tenant
	bob := models.Tenant{ ID: 1, Name: "Bob", Location: &models.Room{ Name: "limbo" }, Doing: &sleeping }

	//Initialize clock cycle
	update := make(chan clockcycle.ClockTime)
	location, _ := time.LoadLocation("EST")
	startTime := time.Date(time.Now().Year(),time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)
	clock := clockcycle.ClockCycle{ StartTime: startTime, Interval: time.Duration(1*time.Second), Frequency: 24*4, Update: update }
	go clock.Start()
	for {
		fmt.Println("-----------------")
		clocktime := <-clock.Update
		interval := clocktime.Interval
		cycle := clocktime.Cycle
		fmt.Println(cycle, interval)
		hour := int(interval / 4)
		minute := (interval % 4) * 15
		fmt.Println(hour, ":",minute)

		if interval == 40 {
			bob.AddAction(interval, 24*4, &goToWork);
			bob.AddAction(interval, 24*4, &goToSleep);
		}
		bob.DoNextAction(interval)
		fmt.Println(bob.Doing, bob.Location)
		bob.PrintSchedule()
	}

}


