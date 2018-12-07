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

func (bprouter *BPRouter) tenantRoutes() {
	//Tenant routes
	tenantRouter := bprouter.Router.PathPrefix("/tenant").Subrouter()
	tenantRouter.HandleFunc("/create/{house_id}/", bprouter.createTenant).Methods("POST", "OPTIONS")
}

func (bprouter *BPRouter) createTenant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//Find house to insert tenant in from url parameters
	houseid, err := strconv.Atoi(params["house_id"])
	if err != nil {
		panic(err)
	}
	house := bprouter.World.Houses[houseid]

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
