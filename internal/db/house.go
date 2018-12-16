package db

import (
	"fmt"
)

//House represents a house
type House struct {
	tableName struct{} `sql:"houses"`
	HouseID   int      `sql:"house_id,type:serial,pk"`
	Name      string   `sql:"name"`
	Width     int      `sql:"width,type:smallint"`
	Height    int      `sql:"height,type:smallint"`
}

//CreateHousesTable creates a houses table in a database
func (bpdb *BitpartmentDB) CreateHousesTable() error {
	return bpdb.createSomeTable(&House{}, "HOUSE")
}

//DropHousesTable will drop the tenants table in the database
func (bpdb *BitpartmentDB) DropHousesTable() error {
	return bpdb.dropSomeTable(&House{}, "HOUSE")
}

//InsertHouse inserts a house into the houses table
func (bpdb *BitpartmentDB) InsertHouse(house *House) error {
	return bpdb.insert(house, "HOUSE")
}

//GetHouseByID returns a House instance from the Houses table based on id
func (bpdb *BitpartmentDB) GetHouseByID(id int) (*House, error) {
	house := &House{HouseID: id}
	err := bpdb.db.Model(house).Where("house_id = ?house_id").Select()
	if err != nil {
		return nil, err
	}
	fmt.Println("Got house:", house)
	return house, nil
}
