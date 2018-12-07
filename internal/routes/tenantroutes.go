package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (server *Server) tenantRoutes() {
	//Tenant routes
	tenantRouter := server.Router.PathPrefix("/tenant").Subrouter()
	tenantRouter.HandleFunc("/create/{house_id}/", server.createTenant).Methods("POST", "OPTIONS")
}

func (server *Server) createTenant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//Find house to insert tenant in from url parameters
	houseid, err := strconv.Atoi(params["house_id"])
	if err != nil {
		panic(err)
	}
	house := server.World.Houses[houseid]

	//Get tenant creation data from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	var tenant models.Tenant
	if err := json.Unmarshal(body, &tenant); err != nil {
		panic(err)
	}
	house.AddTenant(tenant)
}
