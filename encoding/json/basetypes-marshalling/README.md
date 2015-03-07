# encoding/json parsing basetypes

This playbook shows how to parse basetypes with GO using different objects.

We will use the following json:

```
{
  "number": 123,
  "bool": true,
  "string": "hello gopher"
}
```

## Private fields

Private fields are not marshalled/unmarshalled even if annotated. The only way workaround would be to write a custom marshal mehtod for the object and handle it internally.

The following struct:
```
type privateFields struct {
	privateField          int
	privateAnnotatedField int `json:"number"`
}
```
Is marshalled as:
```
&main.privateFields{privateField:0, privateAnnotatedField:0}
```

## Fields annotation

Field annotation is usefult to decouple fields names from their json representation. In most cases that is necessary because GO public fields need to start with capital letter while json fields are written lowercase with the `_` notation.

The following struct:
```
type matchingStruct struct {
	Number int    `json:"number"`
	Bool   bool   `json:"bool"`
	String string `json:"string"`
}
```
Is marshalled as:
```
&main.matchingStruct{Number:123, Bool:true, String:"hello gopher"}
```

## Default values and pointers

Basetypes are initialized even if missing from the json, unless they are pointers.

The following struct:
```
type missingFields struct {
	StringByRef  *string `json:"string"`
	MissingValue string  `json:"another-string-val" omitempty`
	MissingRef   *string `json:"another-string-ref" omitempty`

	BoolByRef          *bool `json:"bool"`
	MissingBoolByValue bool  `json:"missing-bool-val" omitempty`
	MissingBoolByRef   *bool `json:"missing-bool-ref" omitempty`
}
```
Is marshalled as:
```
&main.missingFields{StringByRef:(*string)(0x2081ec4a0), MissingValue:"", MissingRef:(*string)(nil), BoolByRef:(*bool)(0x2081ec496), MissingBoolByValue:false, MissingBoolByRef:(*bool)(nil)}
```

## Duplicated json fields are ingored

Duplicated annotation are not marshalled/unmarshalled:

The following struct:
```
type duplicatedAnnotations struct {
	Number1 int `json:"number"`
	Number2 int `json:"number"`
}
```
Is marshalled as:
```
&main.duplicatedAnnotations{Number1:0, Number2:0}
```
Notice that neither of the fields annotated as `number` has been initialized.