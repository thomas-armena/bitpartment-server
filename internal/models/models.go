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

//House  is a struct that represents a house. 
type House struct {
	ID			int
	Tenants		[]Tenant
}

//Room is a struct that represents a room.
type Room struct{
	Name		string
}

//Doing is a struct that represents an action
type Doing struct{
	Name		string
}

//Action is a struct that represents a change of state for a tenant after an action is performed
type Action struct{
	Interval	int
	Location	*Room
	Doing		*Doing
	Next		*Action
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
//given interval matches the interval the next action performs at.
func (tenant *Tenant) DoNextAction(interval int){
	if tenant.NextAction == nil{
		return
	}
	for tenant.NextAction.Interval == interval {
		tenant.DoAction(*tenant.NextAction)
		tenant.NextAction = tenant.NextAction.Next
		if tenant.NextAction == nil{
			return
		}
	}
}

//AddAction is a function that adds an action to the tenant's schedule
func (tenant *Tenant) AddAction(interval int, frequency int, action *Action){
	if tenant.NextAction == nil {
		tenant.NextAction = action
		return
	}
	var prevAction *Action
	var tempAction = tenant.NextAction
	var actionInterval, tempInterval int
	if action.Interval < interval {
		actionInterval = action.Interval + frequency
	} else {
		actionInterval = action.Interval
	}
	if tempAction.Interval < interval {
		tempInterval = tempAction.Interval + frequency
	} else {
		tempInterval = tempAction.Interval
	}

	for actionInterval > tempInterval {
		prevAction = tempAction
		tempAction = tempAction.Next
		if tempAction == nil {
			prevAction.Next = action
			return
		}
	}
	action.Next = tempAction
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
		fmt.Println(temp.Location, temp.Doing, temp.Interval)
		temp = temp.Next
	}
}
