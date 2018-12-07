package routes

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func (server *Server) websocketRoutes() {
	websocketRouter := server.Router.PathPrefix("/ws").Subrouter()
	websocketRouter.HandleFunc("/connect", server.handleWebSocket)
}

var upgrader = websocket.Upgrader{}

func (server *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("websocket listening ...")
	go func(conn *websocket.Conn) {
		for {
			mType, msg, _ := conn.ReadMessage()
			conn.WriteMessage(mType, msg)
		}
	}(conn)
}
