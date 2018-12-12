package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func (server *Server) websocketRoutes() {
	websocketRouter := server.Router.PathPrefix("/ws").Subrouter()
	websocketRouter.HandleFunc("/connect", server.handleWebSocket)
}

var upgrader = websocket.Upgrader{}

type message struct {
	Username string
	HouseID  int
}

func (server *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	connected := false
	var msg message

	//Initialize connection with world
	go func(server *Server, conn *websocket.Conn, connected *bool, msg *message) {
		mType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(mType, string(data), err)
		if err := json.Unmarshal(data, msg); err != nil {
			panic(err)
		}
		fmt.Println(server.World.OpenConnection((*msg).Username, (*msg).HouseID))
		fmt.Println("connected!")
		*connected = true
		house := server.World.Houses[(*msg).HouseID]
		websocket.WriteJSON(conn, house)
	}(server, conn, &connected, &msg)

	//Wait for a connection to be established
	for {
		if connected {
			break
		}
	}

	//Change connection to world
	go func(server *Server, conn *websocket.Conn, connected *bool, msg *message) {
		for {
			mType, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(mType, string(data), err)
			if err := json.Unmarshal(data, msg); err != nil {
				panic(err)
			}
			fmt.Println(server.World.ChangeConnection((*msg).Username, (*msg).HouseID))
			fmt.Println("connection changed")
			house := server.World.Houses[(*msg).HouseID]
			websocket.WriteJSON(conn, house)
		}
	}(server, conn, &connected, &msg)

	//Start streaming updated data of requested house to client
	go func(server *Server, conn *websocket.Conn, msg *message) {
		for {
			house := <-server.World.Connections[msg.Username].Channel
			fmt.Println(house)
			fmt.Println(server.World.Houses)
			websocket.WriteJSON(conn, house)
		}
	}(server, conn, &msg)
}
