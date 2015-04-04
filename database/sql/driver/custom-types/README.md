#Saving custom types with sql

In golang it is relatively common to "override" base types with custom redefined named types. This is a common practice in those situations where in other languages enumerative would be used, or when trying to model a finite state machine.

By default the sql driver only handle the following set of types.
* int64
* float64
* bool
* []byte
* string
* time.Time
* nil - for NULL values
Attempting to save any other type will return in an error, unless the type implements the right interfaces.
This playbook shows how to do that by storing a custom string to sqlite3.

## Defining a custom named basetype.

In this example we want to define a `Request`, which is an entity with a `code` and a `status`, where the status can only have a fixed set of value. We design that by defining a string type called `RequestStatus`. 

```
type RequestStatus string

const (
	QUEUED   RequestStatus = "queued"
	APPROVED               = "approved"
	REJECTED               = "rejected"
)

var _ALL_STATES = []RequestStatus{QUEUED, APPROVED, REJECTED}

type Request struct {
	Code   string
	Status RequestStatus
}
```

Now in the `Request` object the value of the status field can only have one of the predefined values. This is conceptually similar to an enumerative. Unfortunately if we attempt to store a request object in a SQL databse the database driver will throw an exception comlaining that it doesn't know how to handle variables of type `RequestStatus`.


# Implementing the driver.Valuer interface.

The [Valuer](http://golang.org/pkg/database/sql/driver/#Valuer) interface is used to pass the instance value to the driver. In this example the value is a string, therefore a simple cast would do. 

```
func (s RequestStatus) Value() (driver.Value, error) {
	return string(s), nil
}
```

# Implementing the sql.Scanner interface.

The [Scanner](http://golang.org/pkg/database/sql/#Scanner) interface is used to load from the sql driver. 

```
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
```

Notice that we used reflection to reassign the current value.