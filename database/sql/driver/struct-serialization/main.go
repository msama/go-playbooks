package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"

	"code.google.com/p/go-uuid/uuid"
	"github.com/mattn/go-sqlite3"
)

// The entity that we want to store in a sql database.
//
// Note that Postgres 9.4 supports json and jsonb fields therefore
// it would be possible to define the body field as json and query it
// as part of a json document.
//
// -- begin
// CREATE TABLE entities (
//	  uuid  uuid PRIMARY KEY,
//	  body jsonb DEFAULT '{}'
// );
// -- end
type Entity struct {
	Uuid string      `json:"uuid"`
	Body *EntityBody `json:"body"`
}

type EntityBody struct {
	Address  string `json:"address,omitempty"`
	PostCode string `json:"postcode,omitempty"`
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	State    string `json:"state,omitempty"`
}

// Instructs the driver to threat EntityBody as a []byte.
// Internally encodes it in json.
func (s *EntityBody) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan and EntityBody from a string.
func (s *EntityBody) Scan(value interface{}) error {
	var byteVal []byte
	switch v := value.(type) {
	case string:
		byteVal = []byte(v)
	case []byte:
		byteVal = v
	default:
		return fmt.Errorf("Cannot convert %s to EntityBody. Unrecognised type.", value)
	}

	if err := json.Unmarshal(byteVal, s); err != nil {
		return err
	}
	return nil
}

// The database driver that we will use
var DB_DRIVER string

func init() {
	// We will use sqlite3 as database
	sql.Register(DB_DRIVER, &sqlite3.SQLiteDriver{})
}

const (
	CREATE = "CREATE TABLE IF NOT EXISTS entities (uuid varchar(255) PRIMARY KEY, body varchar(255))"
	INSERT = "INSERT INTO entities (uuid, body) values (?, ?)"
	SELECT = "SELECT body FROM entities WHERE uuid = ?"
)

// Creates the table.
func createTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(CREATE)
}

// Insert a new Entity.
// Notice that the EntityBody is directly sent to the driver.
func insert(db *sql.DB, e *Entity) (sql.Result, error) {
	return db.Exec(INSERT, e.Uuid, e.Body)
}

// Load an Entity.
func queryEntity(db *sql.DB, id string) (*Entity, error) {
	e := &Entity{
		Uuid: id,
	}
	if err := db.QueryRow(SELECT, id).Scan(&e.Body); err != nil {
		return nil, err
	}
	return e, nil
}

// Query the database raw value so we can see what was stored.
func queryRaw(db *sql.DB, id string) (string, error) {
	var r string
	if err := db.QueryRow(SELECT, id).Scan(&r); err != nil {
		return "", err
	}
	return r, nil
}

func main() {
	database, err := sql.Open(DB_DRIVER, "database_sql_driver_struct_serialization.sqlite3")
	if err != nil {
		log.Fatalf("Failed to create db handler: %v.", err)
	}
	defer database.Close()

	if _, err := createTable(database); err != nil {
		log.Fatalf("Failed to create table: %v.", err)
	}

	newEntity := &Entity{
		Uuid: uuid.New(),
		Body: &EntityBody{
			Address:  "21 foo road",
			PostCode: "12345",
			City:     "London",
			Country:  "GB",
		},
	}

	if _, err := insert(database, newEntity); err != nil {
		log.Fatalf("Failed to store: %v.", err)
	}

	rawString, err := queryRaw(database, newEntity.Uuid)
	if err != nil {
		log.Fatalf("Failed to query: %v.", err)
	}
	log.Printf("Raw: %#v", rawString)

	loadedEntity, err := queryEntity(database, newEntity.Uuid)
	if err != nil {
		log.Fatalf("Failed to query: %v.", err)
	}
	log.Printf("Entity: %#v", loadedEntity.Body)
}
