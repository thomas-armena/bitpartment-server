package main

import (
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"github.com/thomas-armena/bitpartment-server/internal/server"
	"github.com/thomas-armena/bitpartment-server/pkg/clockcycle"
	"time"
)

func main() {

	//Initialise database
	bitpartmentDB := db.BitpartmentDB{}
	bitpartmentDB.Connect()
	bitpartmentDB.DropTenantsTable()
	bitpartmentDB.DropHousesTable()
	bitpartmentDB.DropRoomsTable()
	bitpartmentDB.DropActionsTable()
	bitpartmentDB.CreateTenantsTable()
	bitpartmentDB.CreateHousesTable()
	bitpartmentDB.CreateRoomsTable()
	bitpartmentDB.CreateActionsTable()

	//Test queries

	insertBaseHouse(&bitpartmentDB, "basehouse1")
	insertBaseHouse(&bitpartmentDB, "basehouse2")
	insertBaseHouse(&bitpartmentDB, "basehouse3")
	insertBaseHouse(&bitpartmentDB, "basehouse4")
	insertBaseHouse(&bitpartmentDB, "basehouse5")

	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Kristie",
		RoomID:   2,
		HouseID:  1,
		ActionID: -1,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Donald",
		RoomID:   3,
		HouseID:  2,
		ActionID: -1,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Matt",
		RoomID:   1,
		HouseID:  2,
		ActionID: -1,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Toby",
		RoomID:   1,
		HouseID:  1,
		ActionID: -1,
	})
	bitpartmentDB.InsertTenant(&db.Tenant{
		Name:     "Rod",
		RoomID:   5,
		HouseID:  5,
		ActionID: -1,
	})

	//bitpartmentDB.DeleteTenantByID(2)
	a, _ := bitpartmentDB.GetAvailableActionsByHouseID(1)
	fmt.Println(a)

	//Initialize clock cycle
	update := make(chan clockcycle.ClockTime)
	location, _ := time.LoadLocation("EST")
	frq := 24
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)
	clock := clockcycle.ClockCycle{StartTime: startTime, IntervalDuration: time.Duration(1 * time.Second), Frequency: frq, Update: update}
	go clock.Start()

	//Initialize server
	bpserver := server.NewServer(&bitpartmentDB, &clock)
	go bpserver.Run()
	for {
		clocktime := <-clock.Update
		bpserver.UpdateStates(clocktime.Cycle, clocktime.Interval)
		bpserver.DispatchConnections()
		fmt.Println(clocktime)
	}

}

func insertBaseHouse(bpdb *db.BitpartmentDB, name string) {
	model, _ := bpdb.InsertHouse(&db.House{
		Name:   name,
		Width:  50,
		Height: 50,
	})
	house := model.(*db.House)
	insertBaseRoom(bpdb, "gym", house.HouseID)
	insertBaseRoom(bpdb, "living room", house.HouseID)
	insertBaseRoom(bpdb, "bar", house.HouseID)

	fmt.Println(house.HouseID)
	//bpdb.InsertRoom(HouseID)
}

//ActionsInRoom Contains a map of starting actions inside a room
var ActionsInRoom = map[string][]db.Action{
	"gym":         {db.Action{Type: "weight training", Intervals: 1}, db.Action{Type: "cardio training", Intervals: 2}},
	"bar":         {db.Action{Type: "bartending", Intervals: 4}, db.Action{Type: "socializing", Intervals: 3}},
	"living room": {db.Action{Type: "watching tv", Intervals: 2}, db.Action{Type: "socializing", Intervals: 4}},
}

func insertBaseRoom(bpdb *db.BitpartmentDB, name string, houseID int) {
	model, _ := bpdb.InsertRoom(&db.Room{
		HouseID: houseID,
		Type:    name,
		Width:   2,
		Height:  1,
		X:       0,
		Y:       0,
	})
	room := model.(*db.Room)
	roomID := room.RoomID

	for _, action := range ActionsInRoom[name] {
		action.RoomID = roomID
		action.HouseID = houseID
		action.TenantID = -1
		bpdb.InsertAction(&action)
	}

}
