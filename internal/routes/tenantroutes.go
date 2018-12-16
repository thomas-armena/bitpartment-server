package routes

import (
	"encoding/json"
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"io/ioutil"
	"net/http"
)

func (server *Server) tenantRoutes() {
	//Tenant routes
	tenantRouter := server.Router.PathPrefix("/tenant").Subrouter()
	tenantRouter.HandleFunc("/create", server.createTenant).Methods("POST", "OPTIONS")
}

func (server *Server) createTenant(w http.ResponseWriter, r *http.Request) {

	//Get tenant creation data from request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	var tenant db.Tenant
	if err := json.Unmarshal(body, &tenant); err != nil {
		panic(err)
	}
	server.DB.InsertTenant(&tenant)
	respondWithJSON(w, http.StatusCreated, tenant)

}
