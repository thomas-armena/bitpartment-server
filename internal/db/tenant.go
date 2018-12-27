package db

import (
	"fmt"
)

//Tenant is a model for the tenants table
type Tenant struct {
	tableName struct{} `sql:"tenants"`

	TenantID int    `sql:"tenant_id,type:serial,pk"`
	Username string `sql:"username,unique"`
	Name     string `sql:"name"`

	HouseID  int `sql:"house_id"`
	RoomID   int `sql:"room_id,type:smallint"`
	ActionID int `sql:"action_id"`

	NextInterval int `sql:"next_interval"`
	NextDay      int `sql:"next_day"`
}

//CreateTenantsTable creates a tenant table in a database
func (bpdb *BitpartmentDB) CreateTenantsTable() error {
	return bpdb.createSomeTable(&Tenant{}, "TENANT")
}

//DropTenantsTable will drop the tenants table in the database
func (bpdb *BitpartmentDB) DropTenantsTable() error {
	return bpdb.dropSomeTable(&Tenant{}, "TENANT")
}

//InsertTenant inserts a tenant into the tenants table
func (bpdb *BitpartmentDB) InsertTenant(tenant *Tenant) (interface{}, error) {
	return bpdb.insert(tenant, "TENANT")
}

//DeleteTenantByID deletes a tenant that is inside the tenant table
func (bpdb *BitpartmentDB) DeleteTenantByID(id int) error {
	tenant := &Tenant{TenantID: id}
	_, err := bpdb.db.Model(tenant).Where("tenant_id = ?tenant_id").Delete()
	if err != nil {
		return err
	}
	fmt.Println("Deleted tenant")
	return nil
}

//GetTenantByID returns a tenant instance from the Tenants table based on id
func (bpdb *BitpartmentDB) GetTenantByID(id int) (*Tenant, error) {
	tenant := &Tenant{TenantID: id}
	err := bpdb.db.Model(tenant).Where("tenant_id = ?tenant_id").Select()
	if err != nil {
		return nil, err
	}
	fmt.Println("Got tenant:", tenant)
	return tenant, nil
}

//GetTenantsByHouseID returns a list of tenants who live in a given house
func (bpdb *BitpartmentDB) GetTenantsByHouseID(houseID int) ([]Tenant, error) {
	var tenants []Tenant
	err := bpdb.db.Model(&tenants).Where("house_id = ?0", houseID).Select()
	if err != nil {
		return nil, err
	}
	fmt.Println("Got tenants:", tenants)
	return tenants, nil
}

//GetTenants returns all the tenants inside the database
func (bpdb *BitpartmentDB) GetTenants() ([]Tenant, error) {
	var tenants []Tenant
	err := bpdb.db.Model(&tenants).Select()
	if err != nil {
		return nil, err
	}
	fmt.Println("Got tenants:", tenants)
	return tenants, nil
}
