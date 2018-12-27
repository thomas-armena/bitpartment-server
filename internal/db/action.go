package db

//Action represents the outcome and info of an action to be performed
type Action struct {
	tableName struct{} `sql:"actions"`
	ActionID  int      `sql:"action_id,type:serial,pk"`
	RoomID    int      `sql:"room_id"`
	HouseID   int      `sql:"house_id"`
	TenantID  int      `sql:"tenant_id"`
	Type      string   `sql:"type"`
}

//CreateActionsTable creates a room table in a database
func (bpdb *BitpartmentDB) CreateActionsTable() error {
	return bpdb.createSomeTable(&Action{}, "ACTION")
}

//DropActionsTable will drop the room table in the database
func (bpdb *BitpartmentDB) DropActionsTable() error {
	return bpdb.dropSomeTable(&Action{}, "ACTION")
}

//InsertAction inserts a room into the rooms table
func (bpdb *BitpartmentDB) InsertAction(action *Action) (interface{}, error) {
	return bpdb.insert(action, "ACTION")
}

//GetAvailableActionsByHouseID gets actions that are both not occupied
//by a tenant and contain the corresponding HouseID
func (bpdb *BitpartmentDB) GetAvailableActionsByHouseID(houseID int) ([]Action, error) {
	var actions []Action
	err := bpdb.db.Model(&actions).Where("house_id = ?0", houseID).Where("tenant_id = -1").Select()
	if err != nil {
		return nil, err
	}
	return actions, nil
}

//GetActionsByHouseID gets all actions in a house
func (bpdb *BitpartmentDB) GetActionsByHouseID(houseID int) ([]Action, error) {
	var actions []Action
	err := bpdb.db.Model(&actions).Where("house_id = ?0", houseID).Select()
	if err != nil {
		return nil, err
	}
	return actions, nil
}
