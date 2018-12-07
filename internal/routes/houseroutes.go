package routes

import (
	"github.com/gorilla/mux"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"net/http"
	"strconv"
)

func (bprouter *BPRouter) houseRoutes() {

	//House routes
	houseRouter := bprouter.Router.PathPrefix("/house").Subrouter()
	houseRouter.HandleFunc("/create/{id}/", bprouter.createHouse).Methods("GET")
}

func (bprouter *BPRouter) createHouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	houseid, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	house := models.House{ID: houseid}
	bprouter.World.AddHouse(house)

}
