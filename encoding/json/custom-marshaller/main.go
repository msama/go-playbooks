package main

import (
	"encoding/json"
	"fmt"
	"time"
)

var jsonData = []byte(`
{
  "start_date": "2015-03-01"
}
`)

const SHORT_DATE_FORMAT = "2006-01-02"

// Specialised date field
type shortDate struct{ time.Time }

func (s *shortDate) MarshalJSON() ([]byte, error) {
	fmt.Println("shortDate.MarshalJSON()")
	return json.Marshal(s.Format(SHORT_DATE_FORMAT))
}

func (s *shortDate) UnmarshalJSON(b []byte) error {
	fmt.Println("shortDate.UnmarshalJSON()")
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	t, err := time.Parse(SHORT_DATE_FORMAT, raw)
	if err != nil {
		return err
	}
	s.Time = t
	return nil
}

// Exposed object
// Only has fields of standard types
type Example struct {
	StartDate time.Time
}

// Auxiliary specialized object
// This is only necessary to hide the specialised fields.
type auxExample struct {
	StartDate shortDate `json:"start_date"`
}

func (s *Example) MarshalJSON() ([]byte, error) {
	fmt.Println("Example.MarshalJSON()")

	aux := &auxExample{
		StartDate: shortDate{Time: s.StartDate},
	}
	return json.Marshal(aux)
}

func (s *Example) UnmarshalJSON(b []byte) error {
	fmt.Println("Example.UnmarshalJSON()")

	var raw *auxExample
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	s.StartDate = raw.StartDate.Time
	return nil
}

func main() {
	var ex *Example
	if err := json.Unmarshal(jsonData, &ex); err != nil {
		panic(err)
	}
	fmt.Printf("Fields are properly unmarshalled:\n%+v\n", ex)

	val, err := json.Marshal(ex)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Fields are properly marshalled:\n%s\n", val)
}
