package routes

import (
	"fmt"
)

//Connection represents a connection from a client to a house
type Connection struct {
	HouseID      int
	ConnectionID string
	Channel      chan int
}

//DispatchConnections sends house data to all the channels subscribed to websockets
func (server *Server) DispatchConnections() {
	for _, connection := range server.Connections {
		fmt.Println("dispatching", connection)
		connection.Channel <- connection.HouseID
	}
}

//OpenConnection adds a connection to the server
func (server *Server) OpenConnection(connectionID string, houseID int) *Connection {
	fmt.Println("connid", connectionID, "houseID", houseID)
	server.Connections[connectionID] = &Connection{HouseID: houseID, ConnectionID: connectionID, Channel: make(chan int)}
	return server.Connections[connectionID]
}

//ChangeConnection changes the houseID of a connection
func (server *Server) ChangeConnection(connectionID string, houseID int) *Connection {
	server.Connections[connectionID].HouseID = houseID
	fmt.Println("changed connection", connectionID, "to house", houseID)
	return server.Connections[connectionID]
}

//CloseConnection removes a connection from the server
func (server *Server) CloseConnection(connectionID string) {
	_, ok := server.Connections[connectionID]
	if ok {
		delete(server.Connections, connectionID)
	}
}
