package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"log"
	"net/http"
)

//BPRouter is a struct that stores the main route for the server
type BPRouter struct {
	World  *models.World
	Router *mux.Router
}

//NewBPRouter is a constructor for the BPRouter struct
func NewBPRouter(world *models.World) *BPRouter {
	router := mux.NewRouter()
	bprouter := &BPRouter{Router: router, World: world}
	bprouter.houseRoutes()
	bprouter.tenantRoutes()
	bprouter.Router = router
	return bprouter
}

//Run is a function that serves the router
func (bprouter *BPRouter) Run() {
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"http://localhost:3000"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})
	handler := c.Handler(bprouter.Router)
	log.Fatal(http.ListenAndServe(":8000", handler))
	fmt.Println("Listening on port 8000")
}
