# json2args

## Example
### Command
```bash
json2go --package main --typename Person --varname PersonList --output person.go persons.json
```

### Input JSON File
```json
[
    {
        "id": "3a98ee9c-6b60-4239-9d12-cbd637a904c9",
        "name": "alice",
        "age": 33
    },
    {
        "id": "cf17db3d-0cae-44c7-91a1-d834785a25e6",
        "name": "bob",
        "age": 19
    }
]
```

### Output Go File
```go
package main

type Person struct {
	Age  float64
	Id   string
	Name string
}

var PersonList = []Person{
	{
		Age:  33,
		Id:   "3a98ee9c-6b60-4239-9d12-cbd637a904c9",
		Name: "alice",
	},
	{
		Age:  19,
		Id:   "cf17db3d-0cae-44c7-91a1-d834785a25e6",
		Name: "bob",
	},
}
```
