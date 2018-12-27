package server

import (
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"github.com/thomas-armena/bitpartment-server/internal/utils"
	"math/rand"
)

//UpdateStates is a function that will update actions of all the tenants
func (server *Server) UpdateStates(day, interval int) {
	fmt.Println("----")
	server.DB.GetDB().Model(&db.Tenant{}).ForEach(func(tenant db.Tenant) error {
		dayPassed := tenant.NextDay < day
		dayEqual := tenant.NextDay == day
		intervalPassed := tenant.NextInterval <= interval
		if dayPassed || (dayEqual && intervalPassed) {
			executeAnAction(server, &tenant, day, interval)
			err := server.DB.Update(&tenant)
			utils.PanicIfErr(err)
		}
		fmt.Println(tenant)
		return nil
	})
	fmt.Println("----")
}

func executeAnAction(server *Server, tenant *db.Tenant, day, interval int) {
	action := getAQueuedAction(server, tenant)
	if action == nil {
		action = getAFillerAction(server, tenant)
	}
	execute(tenant, action)

	//set next day and interval
	intervals := 5
	tenant.NextInterval = (interval + intervals) % server.Clock.Frequency
	tenant.NextDay = day + int((interval+intervals)/server.Clock.Frequency)
}

func getAQueuedAction(server *Server, tenant *db.Tenant) *db.Action {
	return nil
}

func getAFillerAction(server *Server, tenant *db.Tenant) *db.Action {
	actions, err := server.DB.GetAvailableActionsByHouseID(tenant.HouseID)
	utils.PanicIfErr(err)
	choice := rand.Intn(len(actions))
	return &actions[choice]
}

func execute(tenant *db.Tenant, action *db.Action) {
	tenant.ActionID = action.ActionID
}
