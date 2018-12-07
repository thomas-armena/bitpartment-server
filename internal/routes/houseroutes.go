package routes

import (
	"github.com/gorilla/mux"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"net/http"
	"strconv"
)

func (server *Server) houseRoutes() {

	//House routes
	houseRouter := server.Router.PathPrefix("/house").Subrouter()
	houseRouter.HandleFunc("/create/{id}/", server.createHouse).Methods("GET")
}

func (server *Server) createHouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	houseid, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	house := models.House{ID: houseid}
	server.World.AddHouse(house)

}
