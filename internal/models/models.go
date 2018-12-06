package models

//World is a struct that holds all the houses on the server
type World struct {
	Houses map[int]*House
}

//House  is a struct that represents a house.
type House struct {
	ID      int
	Tenants []Tenant
}

//Room is a struct that represents a room.
type Room struct {
	Name string
}

//Doing is a struct that represents an action
type Doing struct {
	Name string
}

//Action is a struct that represents a change of state for a tenant after an action is performed
type Action struct {
	Interval int
	Cycle    int
	Location *Room
	Doing    *Doing
	Next     *Action
}

//AddHouse inserts a house into the World struct
func (world *World) AddHouse(house House) {
	world.Houses[house.ID] = &house
}

//AddTenant insterts a tenant into a house
func (house *House) AddTenant(tenant Tenant) {
	house.Tenants = append(house.Tenants, tenant)
}
