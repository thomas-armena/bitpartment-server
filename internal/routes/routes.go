package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/thomas-armena/bitpartment-server/internal/db"
	"github.com/thomas-armena/bitpartment-server/pkg/clockcycle"
	"log"
	"net/http"
)

//Server is a struct that stores the main route for the server
type Server struct {
	DB          *db.BitpartmentDB
	Router      *mux.Router
	Clock       *clockcycle.ClockCycle
	Connections map[string]*Connection
}

//NewServer is a constructor for the Server struct
func NewServer(bpdb *db.BitpartmentDB, clock *clockcycle.ClockCycle) *Server {
	router := mux.NewRouter()
	server := &Server{Router: router, DB: bpdb, Clock: clock, Connections: make(map[string]*Connection)}
	server.houseRoutes()
	server.tenantRoutes()
	server.websocketRoutes()
	server.Router = router
	return server
}

//Run is a function that serves the router
func (server *Server) Run() {
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"http://localhost:3000"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})
	handler := c.Handler(server.Router)
	log.Fatal(http.ListenAndServe(":8000", handler))
	fmt.Println("Listening on port 8000")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
