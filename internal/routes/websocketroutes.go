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

	var msg message

	//Initialize connection with world
	mType, data, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mType, string(data), err)
	if err := json.Unmarshal(data, msg); err != nil {
		panic(err)
	}
	fmt.Println(server.OpenConnection(msg.Username, msg.HouseID))
	fmt.Println("connected!")
	house, err := server.DB.GetHouseByID(msg.HouseID)
	if err != nil {
		panic(err)
	}
	websocket.WriteJSON(conn, house)

	//Change connection to world
	go func(server *Server, conn *websocket.Conn, msg *message) {
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
			fmt.Println(server.ChangeConnection((*msg).Username, (*msg).HouseID))
			fmt.Println("connection changed")
			house, err := server.DB.GetHouseByID((*msg).HouseID)
			if err != nil {
				panic(err)
			}
			websocket.WriteJSON(conn, house)
		}
	}(server, conn, &msg)

	//Start streaming updated data of requested house to client
	go func(server *Server, conn *websocket.Conn, msg *message) {
		for {
			houseID := <-server.Connections[msg.Username].Channel
			websocket.WriteJSON(conn, struct{ id int }{id: houseID})
		}
	}(server, conn, &msg)
}
