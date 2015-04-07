# Encoding structs as json fields


Modern SQL database attempt to fill the gap between SQL and NoSQL by allowing to store json documents as fields. Where that is not possible developers have always the option of saving a json encoded representation of an object as string.

Imagine a sql table defined as follows, where `uuid` is the object id used for consitency and table references and `body` is a json representation of the rest of the object. This is similar to what document stores do with documents.

```
CREATE TABLE entities (
  uuid  uuid PRIMARY KEY,
  body jsonb DEFAULT '{}'
);
```

That can be represented with a GO struct like the following in which body actually contains various other fields:

```
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
```

Ideally to store and load an instance of the Entity struct we would want to use the following code. In the example we simply pass an instance of `EntityBody` to the sql driver and the driver itself will handle the serialization for us by delegating the instance.

```
// Storing
db.Exec(INSERT, e.Uuid, e.Body)

// Loading
db.QueryRow(SELECT, e.Uuid).Scan(&e.Body)
```

## Serializing and deserializing

We want body to be serialized as json so that it can be stored in a jsonb field. To do that we simply have to implement `driver.Valuer` and `sql.Scanner`.

```
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
```

Json is only one of many possible ways of serializing struct. In Postrgres 9.4 it is also possible to query the data inside the json field as follows:

```
select body->'address' from entities;
```