package main

import (
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"github.com/thomas-armena/bitpartment-server/internal/routes"
	"github.com/thomas-armena/bitpartment-server/pkg/clockcycle"
	"time"
)

func main() {

	world := models.World{
		Houses:      make(map[int]*models.House),
		Connections: make(map[string]*models.Connection),
	}
	/*
		//Define rooms
		bedroom := models.Room{Name: "bedroom"}
		work := models.Room{Name: "shop"}

		//Define doings
		sleeping := models.Doing{Name: "sleeping"}
		working := models.Doing{Name: "coding"}

		//Define tenant
		bob := models.Tenant{ID: 1, Name: "Bob", Location: &models.Room{Name: "limbo"}, Doing: &sleeping}
	*/

	//Initialize clock cycle
	update := make(chan clockcycle.ClockTime)
	location, _ := time.LoadLocation("EST")
	frq := 12
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)
	clock := clockcycle.ClockCycle{StartTime: startTime, Interval: time.Duration(1 * time.Second), Frequency: frq, Update: update}
	go clock.Start()

	//Initialize server
	bprouter := routes.NewServer(&world, &clock)
	go bprouter.Run()

	for {
		clocktime := <-clock.Update
		fmt.Println(clocktime.Cycle, clocktime.Interval)
		world.Update(clocktime.Cycle, clocktime.Interval)
		world.DispatchConnections()

	}
	/*
		for {
			fmt.Println("-----------------")
			clocktime := <-clock.Update
			interval := clocktime.Interval
			cycle := clocktime.Cycle
			fmt.Println("cycle:", cycle, "interval:", interval)

			goToWork := models.Action{Cycle: cycle + 1, Interval: 4, Doing: &working, Location: &work}
			goToSleep := models.Action{Cycle: cycle, Interval: 11, Doing: &sleeping, Location: &bedroom}
			if interval == 5 {
				bob.AddAction(cycle, interval, frq, &goToWork)
				bob.AddAction(cycle, interval, frq, &goToSleep)
			}
			bob.DoNextAction(cycle, interval)
			//fmt.Println(bob.Doing, bob.Location)
			for id, house := range world.Houses {
				fmt.Println(id, house.Tenants)
			}
			fmt.Println(world.Houses)
		}
	*/
}
