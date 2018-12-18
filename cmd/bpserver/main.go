package main

import (
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"github.com/thomas-armena/bitpartment-server/internal/routes"
	"github.com/thomas-armena/bitpartment-server/pkg/clockcycle"
	"time"
)

func main() {

	//Initialise database
	bitpartmentDB := db.BitpartmentDB{}
	bitpartmentDB.Connect()
	bitpartmentDB.DropTenantsTable()
	bitpartmentDB.DropHousesTable()
	bitpartmentDB.CreateTenantsTable()
	bitpartmentDB.CreateHousesTable()

	//Test queries

	bitpartmentDB.InsertHouse(&db.House{
		Name:   "House1",
		Width:  50,
		Height: 50,
	})
	bitpartmentDB.InsertHouse(&db.House{
		Name:   "House2",
		Width:  100,
		Height: 100,
	})
	bitpartmentDB.InsertHouse(&db.House{
		Name:   "House3",
		Width:  100,
		Height: 100,
	})
	bitpartmentDB.InsertHouse(&db.House{
		Name:   "House4",
		Width:  100,
		Height: 100,
	})

	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Kristie",
		RoomID:   2,
		HouseID:  1,
		ActionID: 4,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Donald",
		RoomID:   3,
		HouseID:  2,
		ActionID: 2,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Matt",
		RoomID:   1,
		HouseID:  2,
		ActionID: 1,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Toby",
		RoomID:   1,
		HouseID:  1,
		ActionID: 1,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Rod",
		RoomID:   5,
		HouseID:  5,
		ActionID: 1,
	})

	//bitpartmentDB.DeleteTenantByID(2)
	fmt.Println(bitpartmentDB.GetTenantByID(1))
	fmt.Println(bitpartmentDB.GetTenantsByHouseID(2))
	fmt.Println(bitpartmentDB.GetHouses())

	//Initialize clock cycle
	update := make(chan clockcycle.ClockTime)
	location, _ := time.LoadLocation("EST")
	frq := 12
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)
	clock := clockcycle.ClockCycle{StartTime: startTime, Interval: time.Duration(1 * time.Second), Frequency: frq, Update: update}
	go clock.Start()

	//Initialize server
	server := routes.NewServer(&bitpartmentDB, &clock)
	go server.Run()
	for {
		clocktime := <-clock.Update
		server.DispatchConnections()
		fmt.Println(clocktime)
	}

}
