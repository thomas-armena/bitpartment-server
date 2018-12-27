package db

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

//BitpartmentDB is a type that holds an instance of a database
type BitpartmentDB struct {
	db *pg.DB
}

//GetDB returns an instance of the database, to be used to perform operations not provided
//in this package
func (bpdb *BitpartmentDB) GetDB() *pg.DB {
	return bpdb.db
}

//Connect connects the instance to the database
func (bpdb *BitpartmentDB) Connect() error {
	opts := &pg.Options{
		User:     "thomas",
		Password: "032595ab",
		Database: "bitpartment",
		Addr:     "localhost:5432",
	}

	bpdb.db = pg.Connect(opts)
	if bpdb == nil {
		return ErrConn
	}
	fmt.Println("Connected to database")
	return nil
}

//Close will close the instance of the database
func (bpdb *BitpartmentDB) Close() error {
	err := bpdb.db.Close()
	if err != nil {
		return ErrClose
	}
	fmt.Println("Connection to database closed")
	return nil
}

//Update the database based on a provided model
func (bpdb *BitpartmentDB) Update(model interface{}) error {
	err := orm.Update(bpdb.db, model)
	if err != nil {
		return err
	}
	return nil
}

//ErrConn is an error that occurs when connection to database fails
var ErrConn = errors.New("Failed to connect to database")

//ErrClose is an error that occurs when closing a database fails
var ErrClose = errors.New("Failed to close database")

func (bpdb *BitpartmentDB) createSomeTable(model interface{}, name string) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	err := bpdb.db.CreateTable(model, opts)
	if err != nil {
		return err
	}
	fmt.Println("Created", name, "table", model)
	return nil
}

func (bpdb *BitpartmentDB) dropSomeTable(model interface{}, name string) error {

	//Drop the table if it exists
	err := bpdb.db.DropTable(model, &orm.DropTableOptions{IfExists: true})
	if err != nil {
		return err
	}
	fmt.Println("Dropped", name, "table")
	return nil
}

func (bpdb *BitpartmentDB) insert(model interface{}, name string) (interface{}, error) {
	err := bpdb.db.Insert(model)
	if err != nil {
		return nil, err
	}
	fmt.Println("Inserted", model, "into", name)
	return model, nil
}
