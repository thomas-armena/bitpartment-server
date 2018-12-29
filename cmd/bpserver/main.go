package main

import (
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"github.com/thomas-armena/bitpartment-server/internal/server"
	"github.com/thomas-armena/bitpartment-server/internal/utils"
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

	insertTenantWithRoom(&bitpartmentDB, &db.Tenant{
		Name:     "Kristie",
		HouseID:  1,
		ActionID: -1,
	})
	insertTenantWithRoom(&bitpartmentDB, &db.Tenant{
		Name:     "Donald",
		HouseID:  2,
		ActionID: -1,
	})
	insertTenantWithRoom(&bitpartmentDB, &db.Tenant{
		Name:     "Matt",
		HouseID:  2,
		ActionID: -1,
	})
	insertTenantWithRoom(&bitpartmentDB, &db.Tenant{
		Name:     "Toby",
		HouseID:  1,
		ActionID: -1,
	})
	insertTenantWithRoom(&bitpartmentDB, &db.Tenant{
		Name:     "Rod",
		HouseID:  5,
		ActionID: -1,
	})

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
	"gym":         {db.Action{Type: "weight training", Intervals: 1, OwnerID: -1}, db.Action{Type: "cardio training", Intervals: 2, OwnerID: -1}},
	"bar":         {db.Action{Type: "bartending", Intervals: 4, OwnerID: -1}, db.Action{Type: "socializing", Intervals: 3, OwnerID: -1}},
	"living room": {db.Action{Type: "watching tv", Intervals: 2, OwnerID: -1}, db.Action{Type: "socializing", Intervals: 4, OwnerID: -1}},
}

func insertBaseRoom(bpdb *db.BitpartmentDB, name string, houseID int) {
	model, err := bpdb.InsertRoom(&db.Room{
		HouseID: houseID,
		Type:    name,
		Width:   2,
		Height:  1,
		X:       0,
		Y:       0,
	})
	utils.PanicIfErr(err)
	room := model.(*db.Room)
	roomID := room.RoomID

	for _, action := range ActionsInRoom[name] {
		action.RoomID = roomID
		action.HouseID = houseID
		action.TenantID = -1
		bpdb.InsertAction(&action)
	}

}

func insertTenantWithRoom(bpdb *db.BitpartmentDB, tenant *db.Tenant) {
	bpdb.InsertTenant(tenant)
	model, err := bpdb.InsertRoom(&db.Room{
		HouseID: tenant.HouseID,
		Type:    tenant.Name + "'s Room",
		Width:   1,
		Height:  1,
		X:       0,
		Y:       0,
	})
	utils.PanicIfErr(err)
	room := model.(*db.Room)
	bpdb.InsertAction(&db.Action{
		RoomID:    room.RoomID,
		HouseID:   tenant.HouseID,
		TenantID:  -1,
		OwnerID:   tenant.TenantID,
		Type:      "sleeping",
		Intervals: 8,
	})
}
