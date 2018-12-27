package server

import (
	"encoding/json"
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"io/ioutil"
	"net/http"
)

func (server *Server) houseRoutes() {

	//House routes
	houseRouter := server.Router.PathPrefix("/house").Subrouter()
	houseRouter.HandleFunc("/create", server.createHouse).Methods("POST")
	houseRouter.HandleFunc("/get", server.getHouses).Methods("GET")
}

func (server *Server) createHouse(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		panic(err)
	}
	fmt.Println("CREATE HOUSE")
	fmt.Println("Request Body:", string(body))
	var house db.House
	if err := json.Unmarshal(body, &house); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		panic(err)
	}
	fmt.Println("House object:", house)

	server.DB.InsertHouse(&house)
	respondWithJSON(w, http.StatusCreated, house)
}

func (server *Server) getHouses(w http.ResponseWriter, r *http.Request) {
	houses, err := server.DB.GetHouses()
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, http.StatusOK, houses)
}
