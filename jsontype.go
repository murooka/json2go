package main

import (
	"fmt"
	"sort"

	"github.com/iancoleman/strcase"
)

type JSONType struct {
	Empty     bool
	Nullable  bool
	IsBoolean bool
	IsInteger bool
	IsNumber  bool
	IsString  bool
	Array     *JSONType
	Object    map[string]*JSONType
}

func (t *JSONType) ToGoType() string {
	var prefix string
	var typ string

	if t.Nullable {
		prefix = "*"
	}

	if t.IsBoolean {
		typ = "bool"
	} else if t.IsInteger {
		typ = "int"
	} else if t.IsNumber {
		typ = "float64"
	} else if t.IsString {
		typ = "string"
	} else if t.Array != nil {
		typ = fmt.Sprintf(`[]%s`, t.Array.ToGoType())
	} else if t.Object != nil {
		s := fmt.Sprintf("struct{\n")
		keys := make([]string, 0, len(t.Object))
		for key := range t.Object {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			typ := t.Object[key]
			opts := ""
			if typ.Empty {
				opts += ",omitempty"
			}
			s += fmt.Sprintf("%s %s `json:\"%s%s\"`\n", strcase.ToCamel(key), typ.ToGoType(), key, opts)
		}
		s += fmt.Sprintf("}\n")
		typ = s
	}

	if typ == "" {
		panic(fmt.Sprintf("cannot serialize type: %#v", t))
	}

	return prefix + typ
}

func (t *JSONType) Merge(u *JSONType) *JSONType {
	if t == nil {
		return u
	}
	if u == nil {
		return t
	}

	res := &JSONType{}
	res.Empty = t.Empty || u.Empty
	res.Nullable = t.Nullable || u.Nullable
	res.IsBoolean = t.IsBoolean || u.IsBoolean
	res.IsInteger = t.IsInteger || u.IsInteger
	res.IsNumber = t.IsNumber || u.IsNumber
	res.IsString = t.IsString || u.IsString
	res.Array = t.Array.Merge(u.Array)
	if len(t.Object) > 0 || len(u.Object) > 0 {
		res.Object = map[string]*JSONType{}

		keysMap := map[string]struct{}{}
		for key := range t.Object {
			keysMap[key] = struct{}{}
		}
		for key := range u.Object {
			keysMap[key] = struct{}{}
		}

		for key := range keysMap {
			l, lok := t.Object[key]
			r, rok := u.Object[key]
			res.Object[key] = l.Merge(r)
			if !lok || !rok {
				res.Object[key].Empty = true
			}
		}
	}
	return res
}

func detectTypeOfItem(v interface{}) (*JSONType, error) {
	a, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("root value must be array")
	}

	var typ *JSONType = nil
	for _, e := range a {
		typ = typ.Merge(detectType(e))
	}

	return typ, nil
}

func detectType(v interface{}) *JSONType {
	t := &JSONType{}

	switch v := v.(type) {
	case map[string]interface{}:
		t.Object = map[string]*JSONType{}
		for key, val := range v {
			other := detectType(val)
			t.Object[key] = t.Object[key].Merge(other)
		}
	case []interface{}:
		for _, val := range v {
			t.Array = t.Array.Merge(detectType(val))
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
