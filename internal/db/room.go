package db

//Room represents a room, what it can do, and its positional parameters
type Room struct {
	tableName struct{} `sql:"rooms"`
	RoomID    int      `sql:"room_id,type:serial,pk"`
	HouseID   int      `sql:"house_id"`
	OwnerID   int      `sql:"owner_id"`
	Type      string   `sql:"type"`
	Width     int      `sql:"width,type:smallint"`
	Height    int      `sql:"height,type:smallint"`
	X         int      `sql:"x,type:smallint"`
	Y         int      `sql:"y,type:smallint"`
}

//CreateRoomsTable creates a room table in a database
func (bpdb *BitpartmentDB) CreateRoomsTable() error {
	return bpdb.createSomeTable(&Room{}, "ROOM")
}

//DropRoomsTable will drop the room table in the database
func (bpdb *BitpartmentDB) DropRoomsTable() error {
	return bpdb.dropSomeTable(&Room{}, "ROOM")
}

//InsertRooms inserts a room into the rooms table
func (bpdb *BitpartmentDB) InsertRooms(room *Room) error {
	return bpdb.insert(room, "ROOM")
}
