# Custom marshaller

This playbook shows how to write custom marshal and unmarshal methods for your GO struct.

The goal is to parse the following json:
```
{
  "start_date": "2015-03-01"
}
```

## Specialising the parsed field.

Declare a specialised type with a specialised marshalling and unmarshalling logic.

```
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
```

## Wrapping the specialised field from other packages

Delegate an auxiliary object to hide specialised types from othe objects. 

```
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
```

## Further readings
http://talks.golang.org/2015/