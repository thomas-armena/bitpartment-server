package models

import "fmt"

//World is a struct that holds all the houses on the server
type World struct {
	Houses      map[int]*House
	Connections map[string]*Connection
}

//Connection represents a connection from a client to a house
type Connection struct {
	HouseID      int
	ConnectionID string
	Channel      chan *House
}

//House  is a struct that represents a house.
type House struct {
	ID      int
	Name    string
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

//Update is function that makes all tenants in the world perform their actions
func (world *World) Update(cycle, interval int) {
	for _, house := range world.Houses {
		for _, tenant := range house.Tenants {
			tenant.DoNextAction(cycle, interval)
		}
	}
}

//DispatchConnections sends house data to all the channels subscribed to websockets
func (world *World) DispatchConnections() {
	for _, connection := range world.Connections {
		fmt.Println("dispatching", connection)
		connection.Channel <- world.Houses[connection.HouseID]
	}
}

//OpenConnection adds a connection to the world
func (world *World) OpenConnection(connectionID string, houseID int) *Connection {
	fmt.Println("connid", connectionID, "houseID", houseID)
	world.Connections[connectionID] = &Connection{HouseID: houseID, ConnectionID: connectionID, Channel: make(chan *House)}
	return world.Connections[connectionID]
}

//CloseConnection removes a connection from the world
func (world *World) CloseConnection(connectionID string) {
	_, ok := world.Connections[connectionID]
	if ok {
		delete(world.Connections, connectionID)
	}
}
