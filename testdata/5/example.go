// Code generated by "json2go --package main --typename Example --varname Examples --structure map --output example.go example.json"; DO NOT EDIT.

package main

type Example struct{
Age int `json:"age"`
Id string `json:"id"`
Name string `json:"name"`
}


var Examples = map[string]Example{
"alice": {
Age: 33,
Id: "3a98ee9c-6b60-4239-9d12-cbd637a904c9",
Name: "alice",
},
"bob": {
Age: 19,
Id: "cf17db3d-0cae-44c7-91a1-d834785a25e6",
Name: "bob",
},
}