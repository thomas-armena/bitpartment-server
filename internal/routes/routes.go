package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"github.com/thomas-armena/bitpartment-server/pkg/clockcycle"
	"log"
	"net/http"
)

//Server is a struct that stores the main route for the server
type Server struct {
	World  *models.World
	Router *mux.Router
	Clock  *clockcycle.ClockCycle
}

//NewServer is a constructor for the Server struct
func NewServer(world *models.World, clock *clockcycle.ClockCycle) *Server {
	router := mux.NewRouter()
	server := &Server{Router: router, World: world, Clock: clock}
	server.houseRoutes()
	server.tenantRoutes()
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
