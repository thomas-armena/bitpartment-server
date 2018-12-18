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
func (server *Server) OpenConnection(connectionID string, houseID int) (*Connection, error) {

	//return error if connection already exists
	if _, ok := server.Connections[connectionID]; ok {
		return nil, &ErrConnExists{connectionID}
	}

	server.Connections[connectionID] = &Connection{HouseID: houseID, ConnectionID: connectionID, Channel: make(chan int)}

	fmt.Println("Opened Connection:", server.Connections[connectionID])
	return server.Connections[connectionID], nil
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

//ErrConnExists happens when a connection that already exists is attempted to be opened
type ErrConnExists struct{ connectionID string }

func (err *ErrConnExists) Error() string {
	return fmt.Sprintf("The connection with id %s already exists", err.connectionID)
}
