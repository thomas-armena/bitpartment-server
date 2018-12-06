package routes

import(
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"fmt"
	"github.com/thomas-armena/bitpartment-server/internal/models"
	"io/ioutil"
)

//BPRouter is a struct that stores the main route for the server
type BPRouter struct {
	World			*models.World
	Router			*mux.Router
}

//NewBPRouter is a constructor for the BPRouter struct
func NewBPRouter(world *models.World) *BPRouter {
	router := mux.NewRouter()
	bprouter :=  &BPRouter{ Router: router, World: world }
	corsObj:=handlers.AllowedOrigins([]string{"*"})
	router.Use(handlers.CORS(corsObj))

	//House routes
	houseRouter := router.PathPrefix("/house").Subrouter()
	houseRouter.HandleFunc("/create/{id}/", bprouter.createHouse).Methods("GET")

	//Tenant routes
	tenantRouter := router.PathPrefix("/tenant").Subrouter()
	tenantRouter.HandleFunc("/create/{house_id}", bprouter.createTenant).Methods("POST")

	bprouter.Router = router
	return bprouter
}

//Run is a function that serves the router 
func (bprouter *BPRouter) Run() {
	http.ListenAndServe(":8000", bprouter.Router)
	fmt.Println("Listening on port 8000")
}

func (bprouter *BPRouter) createHouse(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	houseid, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	house := models.House{ID:houseid}
	bprouter.World.AddHouse(house)

}

func (bprouter *BPRouter) createTenant(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	houseid, err := strconv.Atoi(params["house_id"])
	house := bprouter.World.Houses[houseid]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	var tenant models.Tenant
	if err := json.Unmarshal(body, &tenant); err != nil{
		panic(err)
	}
	fmt.Println("house",house)
	house.AddTenant(tenant)
	fmt.Println("house",house)
	fmt.Println(tenant)
}



