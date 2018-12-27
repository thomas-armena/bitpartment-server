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

	//set next day and interval
	intervals := action.Intervals
	tenant.NextInterval = (interval + intervals) % server.Clock.Frequency
	tenant.NextDay = day + int((interval+intervals)/server.Clock.Frequency)

	execute(server, tenant, action)
}

func execute(server *Server, tenant *db.Tenant, action *db.Action) {

	//Make previous action point to no one
	if tenant.ActionID != -1 {
		prevAction, err := server.DB.GetActionByTenantID(tenant.TenantID)
		utils.PanicIfErr(err)
		prevAction.TenantID = -1
		err = server.DB.Update(prevAction)
		utils.PanicIfErr(err)
	}

	//Make tenant point to new action
	tenant.ActionID = action.ActionID
	err := server.DB.Update(tenant)
	utils.PanicIfErr(err)

	//Make new action point to tenant
	action.TenantID = tenant.TenantID
	err = server.DB.Update(action)
	utils.PanicIfErr(err)
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
