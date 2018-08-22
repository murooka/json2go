package main

import (
	"fmt"
	"sort"

	"github.com/iancoleman/strcase"
)

type JSONType struct {
	Nullable  bool
	IsBoolean bool
	IsInteger bool
	IsNumber  bool
	IsString  bool
	Array     *JSONType
	Object    map[string]*JSONType
}

func (t *JSONType) ToGoType() string {
	if t.IsBoolean {
		return "bool"
	}

	if t.IsInteger {
		return "int"
	}

	if t.IsNumber {
		return "float64"
	}

	if t.IsString {
		return "string"
	}

	if t.Array != nil {
		return fmt.Sprintf(`[]%s`, t.Array.ToGoType())
	}

	if t.Object != nil {
		s := fmt.Sprintf("struct{\n")
		keys := make([]string, 0, len(t.Object))
		for key := range t.Object {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			typ := t.Object[key]
			s += fmt.Sprintf("%s %s `json:\"%s\"`\n", strcase.ToCamel(key), typ.ToGoType(), key)
		}
		s += fmt.Sprintf("}\n")
		return s
	}

	panic(fmt.Sprintf("cannot serialize type: %#v", t))
}

func (t *JSONType) Merge(u *JSONType) *JSONType {
	if t == nil {
		return u
	}
	if u == nil {
		return t
	}

	res := &JSONType{}
	res.Nullable = t.Nullable || u.Nullable
	res.IsBoolean = t.IsBoolean || u.IsBoolean
	res.IsInteger = t.IsInteger || u.IsInteger
	res.IsNumber = t.IsNumber || u.IsNumber
	res.IsString = t.IsString || u.IsString
	res.Array = t.Array.Merge(u.Array)
	res.Object = map[string]*JSONType{}
	for key, typ := range t.Object {
		res.Object[key] = res.Object[key].Merge(typ)
	}
	for key, typ := range u.Object {
		res.Object[key] = res.Object[key].Merge(typ)
	}
	return res
}

func detectTypeOfItem(v interface{}, fields []string) (*JSONType, error) {
	a, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("root value must be array")
	}

	var typ *JSONType = nil
	for _, e := range a {
		typ = typ.Merge(detectType(e, fields))
	}

	return typ, nil
}

func detectType(v interface{}, fields []string) *JSONType {
	t := &JSONType{}

	switch v := v.(type) {
	case map[string]interface{}:
		t.Object = map[string]*JSONType{}
		var keys []string
		if fields == nil {
			keys = make([]string, 0, len(v))
			for key := range v {
				keys = append(keys, key)
			}
		} else {
			keys = fields
		}
		for _, key := range keys {
			val := v[key]
			t.Object[key] = t.Object[key].Merge(detectType(val, fields))
		}
	case []interface{}:
		for _, val := range v {
			t.Array = t.Array.Merge(detectType(val, fields))
		}
	case string:
		t.IsString = true
	case float64:
		if isInt(v) {
			t.IsInteger = true
		} else {
			t.IsNumber = true
		}
	case bool:
		t.IsBoolean = true
	case nil:
		t.Nullable = true
	}

	return t
}

func isInt(v float64) bool {
	return v == float64(int(v))
}
