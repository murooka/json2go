package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDetectTypeOfItem(t *testing.T) {
	cases := []struct {
		json string
		typ  *JSONType
	}{
		{
			json: `[
				{"prop1": true}
			]`,
			typ: &JSONType{Object: map[string]*JSONType{"prop1": &JSONType{IsBoolean: true}}},
		},
		{
			json: `[
				{"prop1": 1}
			]`,
			typ: &JSONType{Object: map[string]*JSONType{"prop1": &JSONType{IsInteger: true}}},
		},
		{
			json: `[
				{"prop1": 3.14}
			]`,
			typ: &JSONType{Object: map[string]*JSONType{"prop1": &JSONType{IsNumber: true}}},
		},
		{
			json: `[
				{"prop1": "foo"}
			]`,
			typ: &JSONType{Object: map[string]*JSONType{"prop1": &JSONType{IsString: true}}},
		},
		{
			json: `[
				{"prop1": ["foo", "bar", "baz"]}
			]`,
			typ: &JSONType{Object: map[string]*JSONType{"prop1": &JSONType{Array: &JSONType{IsString: true}}}},
		},
	}

	for _, c := range cases {
		var v interface{}
		err := json.NewDecoder(strings.NewReader(c.json)).Decode(&v)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		got, err := detectTypeOfItem(v)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if diff := cmp.Diff(c.typ, got); diff != "" {
			t.Fatalf("diff: %s", diff)
		}
	}
}
