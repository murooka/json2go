// Code generated by "json2go --package main --typename UserAgent --varname UserAgents --root /userAgents --structure map,map --output example.go example1.json example2.json"; DO NOT EDIT.

package main

type UserAgent struct {
	Example string `json:"example"`
}

var UserAgents = map[string]map[string]UserAgent{
	"ios": {
		"safari": {
			Example: "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
		},
		"chrome": {
			Example: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		},
	},
	"windows": {
		"chrome": {
			Example: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		},
	},
}