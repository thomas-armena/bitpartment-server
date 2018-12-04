package models

import "fmt"


//Tenant is a struct that represents a tenant.
type Tenant struct{
	ID			int
	Name		string
	Location	*Room
	Doing		*Doing
	NextAction	*Action
}

//DoAction is a function that changes the state of a tenant using a TenantChange struct
func (tenant *Tenant) DoAction(action Action) {
	if action.Location != nil {
		tenant.Location = action.Location
	}
	if action.Doing != nil {
		tenant.Doing = action.Doing
	}
}

//DoNextAction is a function that will perform the next action in the schedule if the
//given interval matches the cycle and interval the next action performs at.
func (tenant *Tenant) DoNextAction(cycle, interval int){
	if tenant.NextAction == nil{
		return
	}
	for tenant.NextAction.Interval == interval && tenant.NextAction.Cycle == cycle {
		tenant.DoAction(*tenant.NextAction)
		tenant.NextAction = tenant.NextAction.Next
		if tenant.NextAction == nil{
			return
		}
	}
}

//AddAction is a function that adds an action to the tenant's schedule
func (tenant *Tenant) AddAction(cycle int, interval int, frequency int, action *Action){

	//If there are no action in tenant's schedule, just add action to schedule
	if tenant.NextAction == nil {
		tenant.NextAction = action
		return
	}
	var prevAction *Action
	var currAction = tenant.NextAction

	cycleAfter :=  action.Cycle > currAction.Cycle
	cycleEqual := action.Cycle == currAction.Cycle
	intervalAfter := action.Interval > currAction.Interval

	for cycleAfter || (cycleEqual && intervalAfter) {
		prevAction = currAction
		currAction = currAction.Next
		if currAction == nil {
			prevAction.Next = action
			return
		}
		cycleAfter =  action.Cycle > currAction.Cycle
		cycleEqual = action.Cycle == currAction.Cycle
		intervalAfter = action.Interval > currAction.Interval
	}
	action.Next = currAction
	if prevAction != nil{
		prevAction.Next = action
	} else {
		tenant.NextAction = action
	}

}

//PrintSchedule is a function used for debugging that prints out the tenant's current schedule
func (tenant *Tenant) PrintSchedule() {
	temp := tenant.NextAction
	for temp != nil {
		fmt.Println(temp.Location, temp.Doing, temp.Interval, temp.Cycle)
		temp = temp.Next
	}
}
