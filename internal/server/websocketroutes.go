package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"github.com/thomas-armena/bitpartment-server/internal/utils"
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

type output struct {
	House    db.House
	Tenants  []db.Tenant
	Actions  []db.Action
	Rooms    []db.Room
	Cycle    int
	Interval int
}

func (server *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.PanicIfErr(err)

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println(code, text)
		fmt.Println("closing")
		return nil
	})

	var msg message

	//Initialize connection with world

	_, data, err := conn.ReadMessage()
	utils.PanicIfErr(err)

	err = json.Unmarshal(data, &msg)
	utils.PanicIfErr(err)

	_, err = server.OpenConnection(msg.Username, msg.HouseID)
	utils.PanicIfErr(err)

	fmt.Println("connected!")
	out := getOutMessage(server, msg.HouseID)
	websocket.WriteJSON(conn, out)

	//Change connection to world
	go func(server *Server, conn *websocket.Conn, msg *message) {
		for {
			_, data, err := conn.ReadMessage()

			//Close the connection if a closing message was recieved
			if err != nil {
				server.CloseConnection(msg.Username)
				return
			}

			err = json.Unmarshal(data, msg)
			utils.PanicIfErr(err)

			server.ChangeConnection(msg.Username, msg.HouseID)

			out := getOutMessage(server, msg.HouseID)
			websocket.WriteJSON(conn, out)
		}
	}(server, conn, &msg)

	//Start streaming updated data of requested house to client
	go func(server *Server, conn *websocket.Conn, msg *message) {
		for {
			houseID := <-server.Connections[msg.Username].Channel
			out := getOutMessage(server, houseID)
			websocket.WriteJSON(conn, out)
		}
	}(server, conn, &msg)
}

func getOutMessage(server *Server, houseID int) output {
	house, err := server.DB.GetHouseByID(houseID)
	if err != nil {
		panic(err)
	}
	tenants, err := server.DB.GetTenantsByHouseID(houseID)
	if err != nil {
		panic(err)
	}
	actions, err := server.DB.GetActionsByHouseID(houseID)
	if err != nil {
		panic(err)
	}
	rooms, err := server.DB.GetRoomsByHouseID(houseID)
	if err != nil {
		panic(err)
	}

	out := output{
		House:    *house,
		Tenants:  tenants,
		Actions:  actions,
		Rooms:    rooms,
		Cycle:    server.Clock.Cycle,
		Interval: server.Clock.Interval,
	}
	return out
}
