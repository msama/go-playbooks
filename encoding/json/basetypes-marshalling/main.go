package main

import (
	"encoding/json"
	"fmt"
)

var jsonData = []byte(`
{
  "number": 123,
  "bool": true,
  "string": "hello gopher"
}
`)

// private fields are only marshalled/unmarshalled if annotated
type privateFields struct {
	privateField          int
	privateAnnotatedField int `json:"number"`
}

// A struct that perfectly match the json
type matchingStruct struct {
	Number int    `json:"number"`
	Bool   bool   `json:"bool"`
	String string `json:"string"`
}

// Missing basetypes are initialized unless pointers
type missingFields struct {
	StringByRef  *string `json:"string"`
	MissingValue string  `json:"another-string-val" omitempty`
	MissingRef   *string `json:"another-string-ref" omitempty`

	BoolByRef          *bool `json:"bool"`
	MissingBoolByValue bool  `json:"missing-bool-val" omitempty`
	MissingBoolByRef   *bool `json:"missing-bool-ref" omitempty`
}

// Duplicated fields are discarded
type duplicatedAnnotations struct {
	Number1 int `json:"number"`
	Number2 int `json:"number"`
}

func main() {
	fmt.Printf("Welcome to the encoding/json parsing-basetypes Playbook.\n")
	fmt.Printf("In this example we will unmarshal the following json in various struct:\n%s\n", jsonData)

	var pf *privateFields
	if err := json.Unmarshal(jsonData, &pf); err != nil {
		panic(err)
	}
	fmt.Printf("Private fields are not marshalled/unmarshalled even if annotated:\n%#v\n\n", pf)

	var ms *matchingStruct
	if err := json.Unmarshal(jsonData, &ms); err != nil {
		panic(err)
	}
	fmt.Printf("Annotated fields are marshalled/unmarshalled independently from their name:\n%#v\n\n", ms)

	var mf *missingFields
	if err := json.Unmarshal(jsonData, &mf); err != nil {
		panic(err)
	}
	fmt.Printf("Basetypes are initialized even if missing from the json, unless they are pointers:\n%#v\n\n", mf)

	var da *duplicatedAnnotations
	if err := json.Unmarshal(jsonData, &da); err != nil {
		panic(err)
	}
	fmt.Printf("Duplicated annotation are not marshalled/unmarshalled:\n%#v\n\n", da)
}
