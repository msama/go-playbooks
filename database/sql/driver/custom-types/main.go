package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"

	"code.google.com/p/go-uuid/uuid"
	"github.com/mattn/go-sqlite3"
)

type RequestStatus string

const (
	QUEUED   RequestStatus = "queued"
	APPROVED               = "approved"
	REJECTED               = "rejected"
)

var _ALL_STATES = []RequestStatus{QUEUED, APPROVED, REJECTED}

// Implement the driver.Valuer interface so that the value can be sent to the driver.
func (s RequestStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Implement the scanner interface so that the value can be assigned.
func (s *RequestStatus) Scan(value interface{}) error {
	e := reflect.ValueOf(s).Elem()

	var strVal string
	// Convert from one of the driver types to RequestStatus.
	switch v := value.(type) {
	case string:
		strVal = v
	case []byte:
		strVal = string(v)
	case nil:
		return fmt.Errorf("RequestStatus cannot be nil.")
	default:
		return fmt.Errorf("Cannot convert %s to RequestStatus. Unrecognised type.", value)
	}

	// In this example our specialised type limits the values assignable to string.
	// Here we enforce that logic by failing if the string is not a valid RequestStatus.
	var knownState = false
	for _, st := range _ALL_STATES {
		if strVal == string(st) {
			knownState = true
			break
		}
	}
	// The value is not a valid RequestStatus.
	// Implement error recovery here or escalate the error.
	if !knownState {
		return fmt.Errorf("Cannot convert %s to RequestStatus. Unrecognised state.", value)
	}

	// We use reflection to reassign the value of the current instance.
	e.SetString(strVal)
	return nil
}

type Request struct {
	Code   string
	Status RequestStatus
}

// The database driver that we will use
var DB_DRIVER string

func init() {
	// We will use sqlite3 as database
	sql.Register(DB_DRIVER, &sqlite3.SQLiteDriver{})
}

// Creates the table.
func createTable(db *sql.DB) (sql.Result, error) {
	return db.Exec("CREATE TABLE IF NOT EXISTS requests (code varchar(255) PRIMARY KEY, status varchar(255))")
}

// Insert a new Request.
func insert(db *sql.DB, r *Request) (sql.Result, error) {
	return db.Exec("INSERT INTO requests (code, status) values (?, ?)", r.Code, r.Status)
}

// Load a request by code.
func query(db *sql.DB, code string) (*Request, error) {
	r := &Request{
		Code: code,
	}
	//var status string
	err := db.QueryRow("SELECT status FROM requests WHERE code = ?", code).Scan(&r.Status)
	if err != nil {
		return nil, err
	}
	//log.Println(status)
	//r.Status = RequestStatus(status)
	return r, nil
}

func main() {
	database, err := sql.Open(DB_DRIVER, "database_sql_driver_custom_types.sqlite3")
	if err != nil {
		log.Fatalf("Failed to create db handler: %v.", err)
	}
	defer database.Close()

	if _, err := createTable(database); err != nil {
		log.Fatalf("Failed to create table: %v.", err)
	}

	newRequest := &Request{
		Code:   uuid.New(),
		Status: QUEUED,
	}

	if _, err := insert(database, newRequest); err != nil {
		log.Fatalf("Failed to store request: %v.", err)
	}

	loadedRequest, err := query(database, newRequest.Code)
	if err != nil {
		log.Fatalf("Failed to query request: %v.", err)
	}

	log.Printf("Request: %+v", loadedRequest)
}
